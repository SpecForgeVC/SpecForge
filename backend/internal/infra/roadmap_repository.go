package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type roadmapItemRepository struct {
	queries *db.Queries
}

func NewRoadmapItemRepository(queries *db.Queries) app.RoadmapItemRepository {
	return &roadmapItemRepository{queries: queries}
}

func (r *roadmapItemRepository) Get(ctx context.Context, id uuid.UUID) (*domain.RoadmapItem, error) {
	row, err := r.queries.GetRoadmapItem(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.RoadmapItem{
		ID:                  row.ID,
		ProjectID:           row.ProjectID,
		Type:                domain.RoadmapItemType(row.Type),
		Title:               row.Title,
		Description:         row.Description.String,
		BusinessContext:     row.BusinessContext.String,
		TechnicalContext:    row.TechnicalContext.String,
		Priority:            domain.RoadmapItemPriority(row.Priority.RoadmapItemPriority),
		Status:              domain.RoadmapItemStatus(row.Status.RoadmapItemStatus),
		RiskLevel:           domain.RiskLevel(row.RiskLevel.RiskLevel),
		BreakingChange:      row.BreakingChange.Bool,
		RegressionSensitive: row.RegressionSensitive.Bool,
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}, nil
}

func (r *roadmapItemRepository) List(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapItem, error) {
	rows, err := r.queries.ListRoadmapItems(ctx, projectID)
	if err != nil {
		return nil, err
	}
	items := make([]domain.RoadmapItem, len(rows))
	for i, row := range rows {
		items[i] = domain.RoadmapItem{
			ID:                  row.ID,
			ProjectID:           row.ProjectID,
			Type:                domain.RoadmapItemType(row.Type),
			Title:               row.Title,
			Description:         row.Description.String,
			BusinessContext:     row.BusinessContext.String,
			TechnicalContext:    row.TechnicalContext.String,
			Priority:            domain.RoadmapItemPriority(row.Priority.RoadmapItemPriority),
			Status:              domain.RoadmapItemStatus(row.Status.RoadmapItemStatus),
			RiskLevel:           domain.RiskLevel(row.RiskLevel.RiskLevel),
			BreakingChange:      row.BreakingChange.Bool,
			RegressionSensitive: row.RegressionSensitive.Bool,
			CreatedAt:           row.CreatedAt.Time,
			UpdatedAt:           row.UpdatedAt.Time,
		}
	}
	return items, nil
}

func (r *roadmapItemRepository) Create(ctx context.Context, item *domain.RoadmapItem) error {
	row, err := r.queries.CreateRoadmapItem(ctx, db.CreateRoadmapItemParams{
		ProjectID:           item.ProjectID,
		Type:                db.RoadmapItemType(item.Type),
		Title:               item.Title,
		Description:         db.TextToSql(item.Description),
		BusinessContext:     db.TextToSql(item.BusinessContext),
		TechnicalContext:    db.TextToSql(item.TechnicalContext),
		Priority:            db.NullRoadmapItemPriority{RoadmapItemPriority: db.RoadmapItemPriority(item.Priority), Valid: true},
		Status:              db.NullRoadmapItemStatus{RoadmapItemStatus: db.RoadmapItemStatus(item.Status), Valid: true},
		RiskLevel:           db.NullRiskLevel{RiskLevel: db.RiskLevel(item.RiskLevel), Valid: true},
		BreakingChange:      db.BoolToSql(item.BreakingChange),
		RegressionSensitive: db.BoolToSql(item.RegressionSensitive),
	})
	if err != nil {
		return err
	}
	item.ID = row.ID
	if row.CreatedAt.Valid {
		item.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		item.UpdatedAt = row.UpdatedAt.Time
	}
	return nil
}

func (r *roadmapItemRepository) Update(ctx context.Context, item *domain.RoadmapItem) error {
	_, err := r.queries.UpdateRoadmapItem(ctx, db.UpdateRoadmapItemParams{
		ID:               item.ID,
		Title:            item.Title,
		Description:      db.TextToSql(item.Description),
		BusinessContext:  db.TextToSql(item.BusinessContext),
		TechnicalContext: db.TextToSql(item.TechnicalContext),
		Status:           db.NullRoadmapItemStatus{RoadmapItemStatus: db.RoadmapItemStatus(item.Status), Valid: true},
	})
	return err
}

func (r *roadmapItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteRoadmapItem(ctx, id)
}
