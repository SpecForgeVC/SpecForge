package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type requirementRepository struct {
	queries *db.Queries
}

func NewRequirementRepository(queries *db.Queries) app.RequirementRepository {
	return &requirementRepository{queries: queries}
}

func (r *requirementRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Requirement, error) {
	row, err := r.queries.GetRequirement(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapRow(row), nil
}

func (r *requirementRepository) List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.Requirement, error) {
	rows, err := r.queries.ListRequirementsByRoadmapItem(ctx, roadmapItemID)
	if err != nil {
		return nil, err
	}
	reqs := make([]domain.Requirement, len(rows))
	for i, row := range rows {
		reqs[i] = *r.mapRow(row)
	}
	return reqs, nil
}

func (r *requirementRepository) Create(ctx context.Context, req *domain.Requirement) error {
	row, err := r.queries.CreateRequirement(ctx, db.CreateRequirementParams{
		RoadmapItemID:      req.RoadmapItemID,
		Title:              req.Title,
		Description:        db.ToSqlString(req.Description),
		Testable:           db.ToSqlBool(req.Testable),
		AcceptanceCriteria: db.ToSqlString(req.AcceptanceCriteria),
		OrderIndex:         db.ToSqlInt32(int32(req.OrderIndex)),
	})
	if err != nil {
		return err
	}
	req.ID = row.ID
	return nil
}

func (r *requirementRepository) Update(ctx context.Context, req *domain.Requirement) error {
	_, err := r.queries.UpdateRequirement(ctx, db.UpdateRequirementParams{
		ID:                 req.ID,
		Title:              req.Title,
		Description:        db.ToSqlString(req.Description),
		Testable:           db.ToSqlBool(req.Testable),
		AcceptanceCriteria: db.ToSqlString(req.AcceptanceCriteria),
		OrderIndex:         db.ToSqlInt32(int32(req.OrderIndex)),
	})
	return err
}

func (r *requirementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteRequirement(ctx, id)
}

func (r *requirementRepository) mapRow(row db.Requirement) *domain.Requirement {
	return &domain.Requirement{
		ID:                 row.ID,
		RoadmapItemID:      row.RoadmapItemID,
		Title:              row.Title,
		Description:        row.Description.String,
		Testable:           row.Testable.Bool,
		AcceptanceCriteria: row.AcceptanceCriteria.String,
		OrderIndex:         int(row.OrderIndex.Int32),
	}
}
