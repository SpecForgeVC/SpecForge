package mcp

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

// Handlers contains the logic for MCP tools
type Handlers struct {
	repo          app.MCPRepository
	sm            *SnapshotStateMachine
	importService app.ImportService
}

func NewHandlers(repo app.MCPRepository, importService app.ImportService) *Handlers {
	return &Handlers{
		repo:          repo,
		sm:            NewSnapshotStateMachine(),
		importService: importService,
	}
}

func (h *Handlers) CreateSnapshot(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input ToolInputCreateSnapshot
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid project_id"}
	}
	rID, err := uuid.Parse(input.RoadmapItemID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid roadmap_item_id"}
	}

	snapshotID, err := h.repo.CreateSnapshot(ctx, pID, rID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to create snapshot", Data: err.Error()}
	}

	return ToolOutputCreateSnapshot{
		SnapshotID: snapshotID.String(),
		ExtractionRequirements: map[string]interface{}{
			"file_tree":           "full",
			"dependencies":        true,
			"api_routes":          true,
			"database_migrations": true,
			"middleware_stack":    true,
			"test_runner":         true,
			"ci_config":           true,
		},
		RequiredSchema: GetEnvironmentSnapshotSchema(),
		StrictNextStep: h.sm.GetRequiredNextStep(domain.StateInitiated),
	}, nil
}

func (h *Handlers) PostSnapshot(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input ToolInputPostSnapshot
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	sID, err := uuid.Parse(input.SnapshotID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid snapshot_id"}
	}

	// Load snapshot to verify state
	snapshot, err := h.repo.GetSnapshot(ctx, sID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Snapshot not found"}
	}

	if err := h.sm.CanTransition(snapshot.State, domain.StateAnalyzing); err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: err.Error()}
	}

	// Save data and set state to Analyzing
	if err := h.repo.SaveSnapshotData(ctx, sID, input.EnvironmentSnapshot); err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to save snapshot data"}
	}

	// Basic scoring based on completeness of the environment snapshot
	conformance := 0.0
	if len(input.EnvironmentSnapshot.FileTree) > 0 {
		conformance += 0.4
	}
	if len(input.EnvironmentSnapshot.APIRoutes) > 0 {
		conformance += 0.3
	}
	if len(input.EnvironmentSnapshot.Database.Migrations) > 0 {
		conformance += 0.3
	}

	scores := domain.Scores{
		RealityConformance:  conformance,
		DependencyIntegrity: 0.9, // Default until detailed dependency analysis is added
		StructuralAlignment: 0.8,
	}

	var verdict string
	if conformance >= 0.7 {
		verdict = "approved"
	} else if conformance >= 0.4 {
		verdict = "requires_review"
	} else {
		verdict = "failed"
	}

	analysis := domain.SnapshotAnalysis{
		SnapshotID:    sID,
		Scores:        scores,
		Verdict:       verdict,
		DriftDetected: false,
	}

	if err := h.repo.SaveAnalysis(ctx, analysis); err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to save analysis"}
	}

	// Update snapshot state to Completed
	h.repo.UpdateSnapshotState(ctx, sID, domain.StateCompleted)

	return ToolOutputPostSnapshot{
		AnalysisResults:  map[string]interface{}{"drift_detected": false},
		Scores:           scores,
		Verdict:          verdict,
		RequiredNextTool: "none",
	}, nil
}

func (h *Handlers) InitProjectImport(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input domain.ToolInputInitProjectImport
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	output, err := h.importService.InitializeImport(ctx, input)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to initialize import", Data: err.Error()}
	}

	return output, nil
}

func (h *Handlers) SubmitProjectSnapshot(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input domain.ToolInputSubmitProjectSnapshot
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	output, err := h.importService.ProcessSnapshot(ctx, input)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to submit snapshot", Data: err.Error()}
	}

	return output, nil
}

func (h *Handlers) GetImportAlignmentRules(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input domain.ToolInputGetImportAlignmentRules
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid project_id"}
	}

	output, err := h.importService.GetAlignmentRules(ctx, pID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to get alignment rules", Data: err.Error()}
	}

	return output, nil
}

func (h *Handlers) SubmitPostImportSnapshot(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input domain.ToolInputSubmitPostImportSnapshot
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	output, err := h.importService.ProcessPostSnapshot(ctx, input)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to submit post-import snapshot", Data: err.Error()}
	}

	return output, nil
}

func (h *Handlers) FinalizeProjectImport(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input domain.ToolInputFinalizeProjectImport
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	output, err := h.importService.FinalizeImport(ctx, input)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to finalize import", Data: err.Error()}
	}

	return output, nil
}

func (h *Handlers) GetSnapshotStatus(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input struct {
		SnapshotID string `json:"snapshot_id"`
	}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	sID, err := uuid.Parse(input.SnapshotID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid snapshot_id"}
	}

	snapshot, err := h.repo.GetSnapshot(ctx, sID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Snapshot not found"}
	}

	return map[string]interface{}{
		"snapshot_id": snapshot.ID,
		"state":       snapshot.State,
	}, nil
}

func (h *Handlers) ListActiveSnapshots(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var input struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid arguments"}
	}

	pID, err := uuid.Parse(input.ProjectID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid project_id"}
	}

	snapshots, err := h.repo.ListActiveSnapshots(ctx, pID)
	if err != nil {
		return nil, &JSONRPCError{Code: -32603, Message: "Failed to list snapshots"}
	}

	return snapshots, nil
}

func (h *Handlers) Help(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"description": "SpecForge Reality Anchor Engine (RAE) MCP Server",
		"usage_order": []string{"create_snapshot", "post_snapshot"},
		"rules": []string{
			"Snapshots must be initiated before posting data.",
			"Snapshot IDs are unique and versioned per roadmap item.",
		},
		"schemas": map[string]interface{}{
			"EnvironmentSnapshot": GetEnvironmentSnapshotSchema(),
		},
	}, nil
}

func GetEnvironmentSnapshotSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"metadata":         map[string]interface{}{"type": "object"},
			"file_tree":        map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
			"dependencies":     map[string]interface{}{"type": "object"},
			"environment":      map[string]interface{}{"type": "object"},
			"api_routes":       map[string]interface{}{"type": "array"},
			"database":         map[string]interface{}{"type": "object"},
			"middleware_stack": map[string]interface{}{"type": "array"},
			"test_framework":   map[string]interface{}{"type": "object"},
			"ci_config":        map[string]interface{}{"type": "object"},
			"exports_surface":  map[string]interface{}{"type": "array"},
		},
	}
}
