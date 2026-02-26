package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type projectRepository struct {
	queries *db.Queries
}

func NewProjectRepository(queries *db.Queries) app.ProjectRepository {
	return &projectRepository{queries: queries}
}

func (r *projectRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	p, err := r.queries.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Project{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		Name:        p.Name,
		Description: p.Description.String,
		TechStack:   db.SqlToJSON(p.TechStack),
		Settings:    db.SqlToJSON(p.Settings),
		MCPSettings: domain.MCPSettings{
			Enabled:       p.McpEnabled.Bool,
			Port:          int(p.McpPort.Int32),
			BindAddress:   p.McpBindAddress.String,
			TokenRequired: p.McpTokenRequired.Bool,
			Token:         p.McpToken.String,
		},
		RepositoryURL: p.RepositoryUrl.String,
		CreatedAt:     p.CreatedAt.Time,
		UpdatedAt:     p.UpdatedAt.Time,
	}, nil
}

func (r *projectRepository) List(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error) {
	rows, err := r.queries.ListProjects(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = domain.Project{
			ID:          row.ID,
			WorkspaceID: row.WorkspaceID,
			Name:        row.Name,
			Description: row.Description.String,
			TechStack:   db.SqlToJSON(row.TechStack),
			Settings:    db.SqlToJSON(row.Settings),
			MCPSettings: domain.MCPSettings{
				Enabled:       row.McpEnabled.Bool,
				Port:          int(row.McpPort.Int32),
				BindAddress:   row.McpBindAddress.String,
				TokenRequired: row.McpTokenRequired.Bool,
				Token:         row.McpToken.String,
			},
			RepositoryURL: row.RepositoryUrl.String,
			CreatedAt:     row.CreatedAt.Time,
			UpdatedAt:     row.UpdatedAt.Time,
		}
	}
	return projects, nil
}

func (r *projectRepository) Create(ctx context.Context, p *domain.Project) error {
	row, err := r.queries.CreateProject(ctx, db.CreateProjectParams{
		WorkspaceID:      p.WorkspaceID,
		Name:             p.Name,
		Description:      db.TextToSql(p.Description),
		TechStack:        db.JSONToSql(p.TechStack),
		Settings:         db.JSONToSql(p.Settings),
		RepositoryUrl:    db.TextToSql(p.RepositoryURL),
		McpEnabled:       sql.NullBool{Bool: p.MCPSettings.Enabled, Valid: true},
		McpPort:          sql.NullInt32{Int32: int32(p.MCPSettings.Port), Valid: true},
		McpBindAddress:   sql.NullString{String: p.MCPSettings.BindAddress, Valid: p.MCPSettings.BindAddress != ""},
		McpTokenRequired: sql.NullBool{Bool: p.MCPSettings.TokenRequired, Valid: true},
		McpToken:         sql.NullString{String: p.MCPSettings.Token, Valid: p.MCPSettings.Token != ""},
	})
	if err != nil {
		return err
	}
	p.ID = row.ID
	if row.CreatedAt.Valid {
		p.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		p.UpdatedAt = row.UpdatedAt.Time
	}
	return nil
}

func (r *projectRepository) Update(ctx context.Context, p *domain.Project) error {
	_, err := r.queries.UpdateProject(ctx, db.UpdateProjectParams{
		ID:               p.ID,
		Name:             p.Name,
		Description:      db.TextToSql(p.Description),
		TechStack:        db.JSONToSql(p.TechStack),
		Settings:         db.JSONToSql(p.Settings),
		RepositoryUrl:    db.TextToSql(p.RepositoryURL),
		McpEnabled:       sql.NullBool{Bool: p.MCPSettings.Enabled, Valid: true},
		McpPort:          sql.NullInt32{Int32: int32(p.MCPSettings.Port), Valid: true},
		McpBindAddress:   sql.NullString{String: p.MCPSettings.BindAddress, Valid: p.MCPSettings.BindAddress != ""},
		McpTokenRequired: sql.NullBool{Bool: p.MCPSettings.TokenRequired, Valid: true},
		McpToken:         sql.NullString{String: p.MCPSettings.Token, Valid: p.MCPSettings.Token != ""},
	})
	return err
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteProject(ctx, id)
}
