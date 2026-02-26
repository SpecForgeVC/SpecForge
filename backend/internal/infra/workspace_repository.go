package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type workspaceRepository struct {
	queries *db.Queries
}

func NewWorkspaceRepository(queries *db.Queries) app.WorkspaceRepository {
	return &workspaceRepository{queries: queries}
}

func (r *workspaceRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	ws, err := r.queries.GetWorkspace(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Workspace{
		ID:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description.String,
		CreatedAt:   ws.CreatedAt.Time,
		UpdatedAt:   ws.UpdatedAt.Time,
	}, nil
}

func (r *workspaceRepository) List(ctx context.Context) ([]domain.Workspace, error) {
	rows, err := r.queries.ListWorkspaces(ctx)
	if err != nil {
		return nil, err
	}
	workspaces := make([]domain.Workspace, len(rows))
	for i, row := range rows {
		workspaces[i] = domain.Workspace{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description.String,
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
	}
	return workspaces, nil
}

func (r *workspaceRepository) Create(ctx context.Context, ws *domain.Workspace) error {
	row, err := r.queries.CreateWorkspace(ctx, db.CreateWorkspaceParams{
		Name:        ws.Name,
		Description: db.TextToSql(ws.Description),
	})
	if err != nil {
		return err
	}
	ws.ID = row.ID
	if row.CreatedAt.Valid {
		ws.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		ws.UpdatedAt = row.UpdatedAt.Time
	}
	return nil
}

func (r *workspaceRepository) Update(ctx context.Context, ws *domain.Workspace) error {
	_, err := r.queries.UpdateWorkspace(ctx, db.UpdateWorkspaceParams{
		ID:          ws.ID,
		Name:        ws.Name,
		Description: db.TextToSql(ws.Description),
	})
	return err
}

func (r *workspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteWorkspace(ctx, id)
}
