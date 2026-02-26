package ui_roadmap

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
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
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
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
