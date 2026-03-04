package drift

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
	"github.com/google/uuid"
)

type ContractRepo interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error)
	List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.ContractDefinition, error)
}

type SnapshotRepo interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.VersionSnapshot, error)
	List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.VersionSnapshot, error)
}

type AuditLogger interface {
	Log(ctx context.Context, entityType string, entityID uuid.UUID, action string, userID uuid.UUID, oldData, newData map[string]interface{}) error
	ListDriftEvents(ctx context.Context) ([]domain.AuditLog, error)
}

type DriftService interface {
	RunDriftCheck(ctx context.Context, contractID uuid.UUID, againstVersionID uuid.UUID) (*domain.DriftReport, error)
	GetFeatureDriftScore(ctx context.Context, featureID uuid.UUID) (int, error)
	GetDriftHistory(ctx context.Context) ([]domain.AuditLog, error)
	GenerateDriftFixes(ctx context.Context, report *domain.DriftReport, roadmapItemID uuid.UUID) ([]domain.DriftFix, error)
}

type driftService struct {
	contractRepo ContractRepo
	snapshotRepo SnapshotRepo
	diffEngine   DiffEngine
	auditLog     AuditLogger
}

func NewDriftService(cRepo ContractRepo, sRepo SnapshotRepo, de DiffEngine, al AuditLogger) DriftService {
	return &driftService{
		contractRepo: cRepo,
		snapshotRepo: sRepo,
		diffEngine:   de,
		auditLog:     al,
	}
}

func (s *driftService) RunDriftCheck(ctx context.Context, contractID uuid.UUID, againstVersionID uuid.UUID) (*domain.DriftReport, error) {
	contract, err := s.contractRepo.Get(ctx, contractID)
	if err != nil {
		return nil, err
	}

	snapshot, err := s.snapshotRepo.Get(ctx, againstVersionID)
	if err != nil {
		return nil, err
	}

	report := &domain.DriftReport{
		DriftDetected:   false,
		BreakingChanges: []domain.BreakingChange{},
		RiskScore:       0.0,
	}

	diffs, err := s.diffEngine.Compare(snapshot.SnapshotData, map[string]interface{}{
		"input":  contract.InputSchema,
		"output": contract.OutputSchema,
	})
	if err != nil {
		return nil, err
	}

	if len(diffs) > 0 {
		report.DriftDetected = true
		for _, d := range diffs {
			report.BreakingChanges = append(report.BreakingChanges, domain.BreakingChange{
				Field: d.Path,
				Issue: d.Description,
			})
			report.RiskScore += float64(d.RiskScore) * 0.1
		}
		if report.RiskScore > 1.0 {
			report.RiskScore = 1.0
		}
		s.auditLog.Log(ctx, "CONTRACT", contractID, "DRIFT_DETECTED", uuid.Nil,
			map[string]interface{}{"version": "current"},
			map[string]interface{}{"drift_report": report},
		)
	}

	return report, nil
}

func (s *driftService) GetFeatureDriftScore(ctx context.Context, featureID uuid.UUID) (int, error) {
	// 1. Get all snapshots, sort newest first
	snapshots, err := s.snapshotRepo.List(ctx, featureID)
	if err != nil {
		return 0, err
	}
	if len(snapshots) == 0 {
		return 100, nil // Fresh feature — no drift baseline yet
	}
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].CreatedAt.After(snapshots[j].CreatedAt)
	})
	latest := snapshots[0]

	// 2. Get current contracts for this feature
	contracts, err := s.contractRepo.List(ctx, featureID)
	if err != nil {
		return 0, err
	}
	if len(contracts) == 0 {
		return 100, nil // No contracts defined yet
	}

	// 3. Extract snapshotted contract state
	snapshotContracts := latest.SnapshotData["contracts"]
	if snapshotContracts == nil {
		return 100, nil // Snapshot predates contract tracking
	}
	snapshotJSON, _ := json.Marshal(snapshotContracts)
	var snapshotMap map[string]interface{}
	if err := json.Unmarshal(snapshotJSON, &snapshotMap); err != nil {
		return 100, nil
	}

	// 4. Compare each current contract against its snapshot state
	totalBreaking, totalNonBreaking := 0, 0
	for _, c := range contracts {
		snapshotState, ok := snapshotMap[c.ID.String()]
		if !ok {
			continue // New contract added since snapshot — not drift
		}
		snapshotStateMap, _ := snapshotState.(map[string]interface{})
		diffs, err := s.diffEngine.Compare(snapshotStateMap, map[string]interface{}{
			"input":  c.InputSchema,
			"output": c.OutputSchema,
		})
		if err != nil {
			continue
		}
		for _, d := range diffs {
			if d.RiskScore >= 3 {
				totalBreaking++
			} else {
				totalNonBreaking++
			}
		}
	}

	// 5. Score: start 100, deduct per change
	score := 100 - (totalBreaking * 15) - (totalNonBreaking * 5)
	if score < 0 {
		score = 0
	}

	// 6. Audit log if drift detected
	if score < 100 {
		s.auditLog.Log(ctx, "FEATURE", featureID, "DRIFT_SCORE_CALCULATED", uuid.Nil,
			nil,
			map[string]interface{}{
				"score":             score,
				"breaking_count":    totalBreaking,
				"nonbreaking_count": totalNonBreaking,
			},
		)
	}

	return score, nil
}

func (s *driftService) GetDriftHistory(ctx context.Context) ([]domain.AuditLog, error) {
	return s.auditLog.ListDriftEvents(ctx)
}

func (s *driftService) GenerateDriftFixes(ctx context.Context, report *domain.DriftReport, roadmapItemID uuid.UUID) ([]domain.DriftFix, error) {
	if report == nil || len(report.BreakingChanges) == 0 {
		return []domain.DriftFix{}, nil
	}

	fixes := make([]domain.DriftFix, 0, len(report.BreakingChanges))
	for _, bc := range report.BreakingChanges {
		fix := domain.DriftFix{
			Field:           bc.Field,
			Issue:           bc.Issue,
			SuggestedChange: "Revert field '" + bc.Field + "' to match the approved contract specification.",
			Explanation:     "The field '" + bc.Field + "' has drifted from the approved spec. Issue: " + bc.Issue + ". Reverting this change will restore contract compliance and prevent downstream failures.",
		}
		fixes = append(fixes, fix)
	}

	s.auditLog.Log(ctx, "ROADMAP_ITEM", roadmapItemID, "DRIFT_FIXES_GENERATED", uuid.Nil,
		nil,
		map[string]interface{}{"fix_count": len(fixes), "risk_score": report.RiskScore},
	)

	return fixes, nil
}
