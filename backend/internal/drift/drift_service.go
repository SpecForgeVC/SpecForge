package drift

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type ContractRepo interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error)
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
	// Get current contract
	contract, err := s.contractRepo.Get(ctx, contractID)
	if err != nil {
		return nil, err
	}

	// Get snapshot for comparison
	snapshot, err := s.snapshotRepo.Get(ctx, againstVersionID)
	if err != nil {
		return nil, err
	}

	// Mocking drift detection logic for now
	// In a real implementation, we would extract the current schema and compare against snapshot data
	report := &domain.DriftReport{
		DriftDetected:   false,
		BreakingChanges: []domain.BreakingChange{},
		RiskScore:       0.0,
	}

	// Compare current contract schema with snapshot
	// (Simplified example)
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

		// Log to Audit
		s.auditLog.Log(ctx, "CONTRACT", contractID, "DRIFT_DETECTED", uuid.Nil,
			map[string]interface{}{"version": "current"},
			map[string]interface{}{"drift_report": report},
		)
	}

	return report, nil
}

func (s *driftService) GetFeatureDriftScore(ctx context.Context, featureID uuid.UUID) (int, error) {
	// 1. Get latest snapshot
	snapshots, err := s.snapshotRepo.List(ctx, featureID)
	if err != nil {
		return 0, err
	}
	if len(snapshots) == 0 {
		return 100, nil // No snapshots means no drift (or fresh feature)
	}
	// latest := snapshots[0] // Assuming sorted desc by date. If not, need sorting.
	// TODO: Ensure repository returns sorted or sort here.

	// 2. Get current contracts (Assuming we want to compare contracts)
	// We need to list contracts for this feature.
	// Problem: ContractRepo interface here only has Get. Need List.
	// For now, let's return 100 to unblock, as we need to update ContractRepo interface and wiring heavily to make this work fully.
	// But let's look at SnapshotData.

	// If SnapshotData has "contracts", we compare.
	// For MVP, if we can't easily fetch current state without circular deps or expanding interfaces too much,
	// we might rely on the last diff report stored?

	// Real implementation:
	// Use s.contractRepo.List(featureID) -> requires updating interface
	// Compare each contract to its version in latest.SnapshotData

	return 100, nil
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

	// Log the fix generation event
	s.auditLog.Log(ctx, "ROADMAP_ITEM", roadmapItemID, "DRIFT_FIXES_GENERATED", uuid.Nil,
		nil,
		map[string]interface{}{"fix_count": len(fixes), "risk_score": report.RiskScore},
	)

	return fixes, nil
}
