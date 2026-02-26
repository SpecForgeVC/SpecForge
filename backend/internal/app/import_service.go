package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/importer"
)

type importService struct {
	projectRepo   ProjectRepository
	mcpRepo       MCPRepository
	alignment     AlignmentService
	diff          DiffService
	bootstrapRepo BootstrapRepository
	sessionRepo   ImportSessionRepository
	scorer        importer.CompletenessScorer
}

func NewImportService(
	projectRepo ProjectRepository,
	mcpRepo MCPRepository,
	alignment AlignmentService,
	diff DiffService,
	bootstrapRepo BootstrapRepository,
	sessionRepo ImportSessionRepository,
) ImportService {
	return &importService{
		projectRepo:   projectRepo,
		mcpRepo:       mcpRepo,
		alignment:     alignment,
		diff:          diff,
		bootstrapRepo: bootstrapRepo,
		sessionRepo:   sessionRepo,
		scorer:        importer.NewCompletenessScorer(),
	}
}

func (s *importService) InitializeImport(ctx context.Context, input domain.ToolInputInitProjectImport) (*domain.ToolOutputInitProjectImport, error) {
	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}

	_, err = s.projectRepo.Get(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	session := &domain.ImportSession{
		ProjectID:         pID,
		CompletenessScore: 0,
		Status:            domain.ImportStatusPartial,
		IterationCount:    0,
		Locked:            false,
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.ToolOutputInitProjectImport{
		Status: "initialized",
		RequiredDocuments: []string{
			"project_overview", "tech_stack", "modules", "apis", "data_models", "contracts", "risks", "change_sensitivity",
		},
		DocumentSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_overview":   map[string]interface{}{"type": "object"},
				"tech_stack":         map[string]interface{}{"type": "object"},
				"modules":            map[string]interface{}{"type": "array"},
				"apis":               map[string]interface{}{"type": "array"},
				"data_models":        map[string]interface{}{"type": "array"},
				"contracts":          map[string]interface{}{"type": "array"},
				"risks":              map[string]interface{}{"type": "array"},
				"change_sensitivity": map[string]interface{}{"type": "array"},
			},
		},
		ScaffoldInstructions: "Please create local structured JSON files in `.specforge/` (e.g., contracts.json, architecture.json, etc.) documenting the requested categories. Follow the iterative IPCP protocol. Read the tool response after submitting a batch, check for missing categories, and submit further batches until the server confirms the project is fully catalogued.",
		SubmissionProtocol: map[string]interface{}{
			"incremental":              true,
			"requires_self_assessment": true,
		},
		SessionID: session.ID.String(),
	}, nil
}

func (s *importService) ProcessSnapshot(ctx context.Context, input domain.ToolInputSubmitProjectSnapshot) (*domain.ToolOutputSubmitProjectSnapshot, error) {
	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}

	session, err := s.sessionRepo.GetLatestSessionByProject(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("no active import session found for project")
	}

	if session.Locked {
		return nil, fmt.Errorf("import session is locked and complete")
	}

	if len(input.SnapshotPayload) == 0 {
		return nil, fmt.Errorf("snapshot_payload is empty: no documents were provided in the request. please ensure you document and submit at least one category (e.g. contracts, apis, modules)")
	}

	// 1. Process Artifact
	artifact := &domain.ImportArtifact{
		SessionID: session.ID,
		Payload:   input.SnapshotPayload,
	}
	if err := s.sessionRepo.CreateArtifact(ctx, artifact); err != nil {
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}

	// 2. Fetch all session artifacts
	artifacts, err := s.sessionRepo.ListArtifactsBySession(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	// 3. Merge payloads
	mergedDocs := make(map[string]interface{})
	for _, a := range artifacts {
		for k, v := range a.Payload {
			if arr, ok := v.([]interface{}); ok {
				if existingArr, exOk := mergedDocs[k].([]interface{}); exOk {
					mergedDocs[k] = append(existingArr, arr...)
				} else {
					mergedDocs[k] = arr
				}
			} else {
				mergedDocs[k] = v
			}
		}
	}

	// 4. Score Completeness
	result := s.scorer.ScoreSubmission(mergedDocs)

	session.IterationCount++
	session.CompletenessScore = int32(result.Score)

	if input.FinalSubmission && result.Score >= 95 {
		session.Locked = true
		session.Status = domain.ImportStatusComplete
	}

	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	catalogState := "draft"
	if session.Locked {
		catalogState = "locked"
	}

	return &domain.ToolOutputSubmitProjectSnapshot{
		Status:                       "accepted",
		CompletenessScore:            result.Score,
		MissingCategories:            result.MissingCategories,
		UnresolvedReferences:         result.UnresolvedReferences,
		SelfAssessmentPrompt:         result.Prompt,
		RequiresAdditionalSubmission: !session.Locked,
		CatalogState:                 catalogState,
		ProjectImportStatus:          string(session.Status),
	}, nil
}

func (s *importService) GetAlignmentRules(ctx context.Context, projectID uuid.UUID) (*domain.ToolOutputGetImportAlignmentRules, error) {
	return &domain.ToolOutputGetImportAlignmentRules{
		StrictRules: []string{
			"No new API routes without contract definition",
			"All database changes must include migration script",
		},
		ForbiddenActions: []string{
			"Direct mutation of global state without variable tracking",
		},
		RequiredSnapshotPolicy: domain.SnapshotPolicy{
			MustCreateSnapshotBeforeChanges: true,
			MustPostSnapshotAfterChanges:    true,
		},
		AlignmentConstraints: map[string]interface{}{
			"layer_isolation":  true,
			"contract_locking": true,
		},
	}, nil
}

func (s *importService) ProcessPostSnapshot(ctx context.Context, input domain.ToolInputSubmitPostImportSnapshot) (*domain.ToolOutputSubmitPostImportSnapshot, error) {
	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}

	// For post-snapshot, we compare against the previous session snapshot
	// or the intended SpecForge state.
	drift, err := s.diff.CompareProjectSnapshot(ctx, input.SnapshotPayload, pID)
	if err != nil {
		return nil, fmt.Errorf("failed to detect post-import drift: %w", err)
	}

	return &domain.ToolOutputSubmitPostImportSnapshot{
		DriftScore:               int(drift.RiskScore),
		ContractBreakageDetected: len(drift.BreakingChanges) > 0,
		RiskDelta:                drift.RiskScore / 10.0,
		AlignmentDelta:           0.0,
		Recommendation:           "approve",
	}, nil
}

func (s *importService) FinalizeImport(ctx context.Context, input domain.ToolInputFinalizeProjectImport) (*domain.ToolOutputFinalizeProjectImport, error) {
	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}

	session, err := s.sessionRepo.GetLatestSessionByProject(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("no active import session found for project")
	}

	// Lock the session and mark as complete
	session.Locked = true
	session.Status = domain.ImportStatusComplete

	if err := s.sessionRepo.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &domain.ToolOutputFinalizeProjectImport{
		Status:      "finalized",
		RedirectURL: fmt.Sprintf("/workspaces/%s/projects/%s/dashboard", "current", pID.String()), // Placeholder dashboard link
	}, nil
}

func getStrFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
