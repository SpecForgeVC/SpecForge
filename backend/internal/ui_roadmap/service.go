package ui_roadmap

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

// Repository defines data access for UI Roadmap Items
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*UIRoadmapItem, error)
	List(ctx context.Context, projectID uuid.UUID) ([]UIRoadmapItem, error)
	Create(ctx context.Context, item *UIRoadmapItem) error
	Update(ctx context.Context, item *UIRoadmapItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Get(ctx context.Context, id uuid.UUID) (*UIRoadmapItem, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, project_id, linked_feature_id, name, description, user_persona, use_case, screen_type, layout_definition, component_tree, state_machine, backend_bindings, accessibility_spec, responsive_spec, validation_rules, animation_rules, design_tokens_used, edge_cases, test_scenarios, intelligence_score, version, created_at, updated_at FROM ui_roadmap_items WHERE id = $1", id)
	var item UIRoadmapItem
	err := row.Scan(&item.ID, &item.ProjectID, &item.LinkedFeatureID, &item.Name, &item.Description, &item.UserPersona, &item.UseCase, &item.ScreenType, &item.LayoutDefinition, &item.ComponentTree, &item.StateMachine, &item.BackendBindings, &item.AccessibilitySpec, &item.ResponsiveSpec, &item.ValidationRules, &item.AnimationRules, &item.DesignTokensUsed, &item.EdgeCases, &item.TestScenarios, &item.IntelligenceScore, &item.Version, &item.CreatedAt, &item.UpdatedAt)
	return &item, err
}

func (r *repo) List(ctx context.Context, projectID uuid.UUID) ([]UIRoadmapItem, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, project_id, linked_feature_id, name, description, user_persona, use_case, screen_type, layout_definition, component_tree, state_machine, backend_bindings, accessibility_spec, responsive_spec, validation_rules, animation_rules, design_tokens_used, edge_cases, test_scenarios, intelligence_score, version, created_at, updated_at FROM ui_roadmap_items WHERE project_id = $1", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UIRoadmapItem
	for rows.Next() {
		var item UIRoadmapItem
		err := rows.Scan(&item.ID, &item.ProjectID, &item.LinkedFeatureID, &item.Name, &item.Description, &item.UserPersona, &item.UseCase, &item.ScreenType, &item.LayoutDefinition, &item.ComponentTree, &item.StateMachine, &item.BackendBindings, &item.AccessibilitySpec, &item.ResponsiveSpec, &item.ValidationRules, &item.AnimationRules, &item.DesignTokensUsed, &item.EdgeCases, &item.TestScenarios, &item.IntelligenceScore, &item.Version, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *repo) Create(ctx context.Context, item *UIRoadmapItem) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ui_roadmap_items (id, project_id, linked_feature_id, name, description, user_persona, use_case, screen_type, layout_definition, component_tree, state_machine, backend_bindings, accessibility_spec, responsive_spec, validation_rules, animation_rules, design_tokens_used, edge_cases, test_scenarios, intelligence_score, version, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`,
		item.ID, item.ProjectID, item.LinkedFeatureID, item.Name, item.Description, item.UserPersona, item.UseCase, item.ScreenType, item.LayoutDefinition, item.ComponentTree, item.StateMachine, item.BackendBindings, item.AccessibilitySpec, item.ResponsiveSpec, item.ValidationRules, item.AnimationRules, item.DesignTokensUsed, item.EdgeCases, item.TestScenarios, item.IntelligenceScore, item.Version, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *repo) Update(ctx context.Context, item *UIRoadmapItem) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ui_roadmap_items SET linked_feature_id=$1, name=$2, description=$3, user_persona=$4, use_case=$5, screen_type=$6, layout_definition=$7, component_tree=$8, state_machine=$9, backend_bindings=$10, accessibility_spec=$11, responsive_spec=$12, validation_rules=$13, animation_rules=$14, design_tokens_used=$15, edge_cases=$16, test_scenarios=$17, intelligence_score=$18, version=$19, updated_at=$20 WHERE id=$21`,
		item.LinkedFeatureID, item.Name, item.Description, item.UserPersona, item.UseCase, item.ScreenType, item.LayoutDefinition, item.ComponentTree, item.StateMachine, item.BackendBindings, item.AccessibilitySpec, item.ResponsiveSpec, item.ValidationRules, item.AnimationRules, item.DesignTokensUsed, item.EdgeCases, item.TestScenarios, item.IntelligenceScore, item.Version, time.Now(), item.ID)
	return err
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM ui_roadmap_items WHERE id = $1", id)
	return err
}

// Service handles business logic for UI-FRE
type Service interface {
	GetItem(ctx context.Context, id uuid.UUID) (*UIRoadmapItem, error)
	ListItems(ctx context.Context, projectID uuid.UUID) ([]UIRoadmapItem, error)
	SaveItem(ctx context.Context, item *UIRoadmapItem) error
	Export(ctx context.Context, id uuid.UUID) (ExportBundle, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	SyncFigma(ctx context.Context, itemID uuid.UUID, payload FigmaSyncPayload) error
	GetPluginAssets(ctx context.Context, itemID uuid.UUID) (map[string]string, error)
	RecommendComponentTree(ctx context.Context, item *UIRoadmapItem) (map[string]interface{}, error)
	RecommendStateMachine(ctx context.Context, item *UIRoadmapItem) (map[string]interface{}, error)
	RecommendFix(ctx context.Context, item *UIRoadmapItem, issues []DriftIssue) (*UIRoadmapItem, error)
	CheckCompliance(ctx context.Context, item *UIRoadmapItem) ([]DriftIssue, error)
	RecommendAPIContracts(ctx context.Context, itemID uuid.UUID) (*domain.APIRecommendation, error)
}

type service struct {
	repo       Repository
	llmService app.LLMService
	rmRepo     app.RoadmapItemRepository
	cRepo      app.ContractRepository
}

func NewService(repo Repository, llm app.LLMService, rmRepo app.RoadmapItemRepository, cRepo app.ContractRepository) Service {
	return &service{
		repo:       repo,
		llmService: llm,
		rmRepo:     rmRepo,
		cRepo:      cRepo,
	}
}

func (s *service) GetItem(ctx context.Context, id uuid.UUID) (*UIRoadmapItem, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) ListItems(ctx context.Context, projectID uuid.UUID) ([]UIRoadmapItem, error) {
	return s.repo.List(ctx, projectID)
}

func (s *service) SaveItem(ctx context.Context, item *UIRoadmapItem) error {
	// 1. Validate
	res := ValidateUIRoadmapItem(item)
	if !res.Valid {
		return fmt.Errorf("validation failed: %v", res.Errors)
	}

	// 2. Score
	scores := CalculateScores(item)
	item.IntelligenceScore = scores.AggregateReadinessScore

	// 3. Save
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
		item.Version = 1
		return s.repo.Create(ctx, item)
	}
	item.UpdatedAt = time.Now()
	return s.repo.Update(ctx, item)
}

func (s *service) Export(ctx context.Context, id uuid.UUID) (ExportBundle, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return ExportBundle{}, err
	}
	return GenerateExportBundle(item)
}

func (s *service) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) SyncFigma(ctx context.Context, itemID uuid.UUID, payload FigmaSyncPayload) error {
	// 1. Fetch the existing item
	item, err := s.repo.Get(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to fetch item for sync: %w", err)
	}

	// 2. Map Hierarchy to ComponentTree (Simplified for this version)
	item.ComponentTree = payload.Hierarchy

	// 3. Update the item
	return s.repo.Update(ctx, item)
}

func (s *service) GetPluginAssets(ctx context.Context, itemID uuid.UUID) (map[string]string, error) {
	manifest := GeneratePluginManifest(itemID.String())
	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")

	code := GeneratePluginCode("http://localhost:8080", itemID.String()) // Base URL should be dynamic in prod

	return map[string]string{
		"manifest.json": string(manifestJSON),
		"code.js":       code,
		"ui.html":       `<div><h3>SpecForge Sync</h3><button onclick="parent.postMessage({pluginMessage: {type: 'sync-hierarchy'}}, '*')">Sync Selected Hierarchy</button></div>`,
	}, nil
}

func (s *service) RecommendComponentTree(ctx context.Context, item *UIRoadmapItem) (map[string]interface{}, error) {
	client, err := s.llmService.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`You are an expert UI Architect. Analyze the following requirements and generate a deterministic, hierarchical component tree.

Identity & Context:
- Name: %s
- Description: %s
- Persona: %s
- Use Case: %s
- Screen Type: %s

Guidelines:
- Return a nested JSON structure where each node has: 'type' (e.g., 'Container', 'Header', 'Sidebar', 'DataGrid', 'Chart', 'ButtonGroup', 'InputGroup', 'Avatar', 'Form', 'FormField'), 'label' (descriptive name), and optional 'children' array.
- The root must be a 'Root' container.
- Ensure the hierarchy is functional and adheres to best practices for the specified screen type (%s).
- Also include a field 'layout_definition' which is a JSON object describing the grid or flex properties (e.g., {"type": "grid", "columns": 12, "gap": "1rem"}).
- Return ONLY a JSON object with two fields: 'component_tree' (containing the root node) and 'layout_definition'. No markdown formatting, NO COMMENTS (// or /*), and strictly valid JSON.`,
		item.Name, item.Description, item.UserPersona, item.UseCase, item.ScreenType, item.ScreenType)

	resp, err := client.Generate(ctx, prompt)
	if err != nil {
		// Fallback for demo
		return map[string]interface{}{
			"type":  "Root",
			"label": "Auto-Generated Layout",
			"children": []interface{}{
				map[string]interface{}{"type": "Header", "label": "Navigation"},
				map[string]interface{}{"type": "Container", "label": "Content Area", "children": []interface{}{
					map[string]interface{}{"type": "DataGrid", "label": "Item List"},
				}},
			},
		}, nil
	}

	cleanedResp := app.CleanJSON(resp)
	var result struct {
		ComponentTree    map[string]interface{} `json:"component_tree"`
		LayoutDefinition map[string]interface{} `json:"layout_definition"`
	}

	if err := json.Unmarshal([]byte(cleanedResp), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Persist layout definition if provided
	if len(result.LayoutDefinition) > 0 {
		layoutJSON, _ := json.Marshal(result.LayoutDefinition)
		item.LayoutDefinition = layoutJSON
	}

	return result.ComponentTree, nil
}

func (s *service) RecommendStateMachine(ctx context.Context, item *UIRoadmapItem) (map[string]interface{}, error) {
	client, err := s.llmService.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`You are an expert UI logic architect. Analyze the following screen requirements and component structure, then generate a comprehensive deterministic state machine.

Identity & Context:
- Name: %s
- Description: %s
- Persona: %s
- Use Case: %s
- Screen Type: %s

Visual Hierarchy (Context for states):
%s

Guidelines:
- Return a JSON object with two fields: 'states' and 'transitions'.
- 'states' is an object where keys are state names. You MUST include exactly these 6 states: 'idle', 'loading', 'success', 'error', 'empty', 'disabled'.
- Each state object contains: 'visual_changes' (string description), 'interaction_changes' (string description), and 'messaging' (string description).
- 'transitions' is an array of objects: { 'from': 'stateA', 'to': 'stateB', 'trigger': 'eventName' }.
- Ensure the state machine handles common UI patterns for the given use case.
- Return ONLY the JSON object. No markdown formatting, NO COMMENTS (// or /*), and strictly valid JSON.`,
		item.Name, item.Description, item.UserPersona, item.UseCase, item.ScreenType, string(item.ComponentTree))

	resp, err := client.Generate(ctx, prompt)
	if err != nil {
		// Fallback for demo with all mandatory states
		return map[string]interface{}{
			"states": map[string]interface{}{
				"idle": map[string]interface{}{
					"visual_changes":      "Default view",
					"interaction_changes": "All enabled",
					"messaging":           "Ready",
				},
				"loading": map[string]interface{}{
					"visual_changes":      "Show skeleton screens",
					"interaction_changes": "Disable inputs",
					"messaging":           "Fetching data...",
				},
				"success": map[string]interface{}{
					"visual_changes":      "Show success indicator",
					"interaction_changes": "Enable actions",
					"messaging":           "Action successful",
				},
				"error": map[string]interface{}{
					"visual_changes":      "Highlight error fields",
					"interaction_changes": "Enable retry",
					"messaging":           "Something went wrong",
				},
				"empty": map[string]interface{}{
					"visual_changes":      "Show empty state illustration",
					"interaction_changes": "Show primary action only",
					"messaging":           "No items found",
				},
				"disabled": map[string]interface{}{
					"visual_changes":      "Opacity reduction",
					"interaction_changes": "Lock all inputs",
					"messaging":           "Feature currently unavailable",
				},
			},
			"transitions": []interface{}{
				map[string]interface{}{"from": "idle", "to": "loading", "trigger": "FETCH_DATA"},
				map[string]interface{}{"from": "loading", "to": "success", "trigger": "DATA_LOADED"},
				map[string]interface{}{"from": "loading", "to": "error", "trigger": "FETCH_FAILED"},
				map[string]interface{}{"from": "loading", "to": "empty", "trigger": "NO_DATA"},
			},
		}, nil
	}

	cleanedResp := app.CleanJSON(resp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleanedResp), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI state machine response: %w", err)
	}

	return result, nil
}
func (s *service) RecommendFix(ctx context.Context, item *UIRoadmapItem, issues []DriftIssue) (*UIRoadmapItem, error) {
	client, err := s.llmService.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	issuesJSON, _ := json.Marshal(issues)
	itemJSON, _ := json.Marshal(item)

	prompt := fmt.Sprintf(`You are a Senior UI/UX Engineer and Logic Architect. A UI Roadmap Item has been identified with several "Drift" or Compliance issues.
Your task is to take the current specification and the detected issues, and provide a RECENT version of the specification that resolves all issues.

Current Specification:
%s

Detected Issues:
%s

Guidelines:
- If a state is missing in the state machine, add it.
- If a component has an invalid binding, fix the binding or the tree structure.
- If accessibility rules are violated, update them.
- Return ONLY the COMPLETE updated JSON object for the entire UIRoadmapItem. No markdown formatting, NO COMMENTS, and strictly valid JSON.`,
		string(itemJSON), string(issuesJSON))

	resp, err := client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	cleanedResp := app.CleanJSON(resp)
	var updatedItem UIRoadmapItem
	if err := json.Unmarshal([]byte(cleanedResp), &updatedItem); err != nil {
		return nil, fmt.Errorf("failed to parse AI fixed response: %w", err)
	}

	return &updatedItem, nil
}

func (s *service) CheckCompliance(ctx context.Context, item *UIRoadmapItem) ([]DriftIssue, error) {
	// 1. Run strict validation
	valRes := ValidateUIRoadmapItem(item)
	var issues []DriftIssue
	for _, errStr := range valRes.Errors {
		issues = append(issues, DriftIssue{
			Type:        "VALIDATION_ERROR",
			Field:       "general",
			Description: errStr,
			Severity:    "ERROR",
		})
	}

	// 2. Run drift detection (placeholder contracts for now)
	driftIssues := DetectUIDrift(item, nil)
	issues = append(issues, driftIssues...)

	return issues, nil
}

func (s *service) RecommendAPIContracts(ctx context.Context, itemID uuid.UUID) (*domain.APIRecommendation, error) {
	// 1. Fetch the UI item
	item, err := s.repo.Get(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// 2. Generate contracts via AI
	client, err := s.llmService.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	itemJSON, _ := json.Marshal(item)
	prompt := fmt.Sprintf(`You are a Senior Systems Architect. Based on the following UI Roadmap Item (Component Tree and State Machine), identify all necessary backend API endpoints required to support this UI feature.

UI Specification:
%s

Guidelines:
- Identify actions that require backend interaction (data fetching, submissions, status checks).
- For each endpoint, provide:
  - Path (e.g., /api/v1/user/profile)
  - Method (GET, POST, etc.)
  - Brief Description of the data contract.
- Return ONLY a JSON object with "title", "description", and "technical_context" (markdown list of endpoints).
- Format:
  {
    "title": string,
    "description": string,
    "technical_context": string (markdown)
  }`, string(itemJSON))

	resp, err := client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	cleanedResp := app.CleanJSON(resp)
	var rec domain.APIRecommendation
	if err := json.Unmarshal([]byte(cleanedResp), &rec); err != nil {
		return nil, fmt.Errorf("failed to parse AI recommendation: %w", err)
	}

	return &rec, nil
}
