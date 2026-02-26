package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type projectService struct {
	repo     ProjectRepository
	auditLog AuditLogService
}

func NewProjectService(repo ProjectRepository, auditLog AuditLogService) ProjectService {
	return &projectService{
		repo:     repo,
		auditLog: auditLog,
	}
}

func (s *projectService) GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return s.repo.Get(ctx, id)
}

func (s *projectService) ListProjects(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error) {
	return s.repo.List(ctx, workspaceID)
}

func (s *projectService) CreateProject(ctx context.Context, workspaceID uuid.UUID, name, description string, techStack map[string]interface{}, settings map[string]interface{}, mcpSettings domain.MCPSettings, repositoryURL string, userID uuid.UUID) (*domain.Project, error) {
	p := &domain.Project{
		WorkspaceID:   workspaceID,
		Name:          name,
		Description:   description,
		TechStack:     techStack,
		Settings:      settings,
		MCPSettings:   mcpSettings,
		RepositoryURL: repositoryURL,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "project", p.ID, "CREATE", userID, nil, map[string]interface{}{"name": p.Name})
	return p, nil
}

func (s *projectService) UpdateProject(ctx context.Context, id uuid.UUID, name, description string, techStack map[string]interface{}, settings map[string]interface{}, mcpSettings domain.MCPSettings, repositoryURL string, userID uuid.UUID) (*domain.Project, error) {
	oldP, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	p := *oldP
	p.Name = name
	p.Description = description
	p.TechStack = techStack
	p.Settings = settings
	p.MCPSettings = mcpSettings
	p.RepositoryURL = repositoryURL
	if err := s.repo.Update(ctx, &p); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "project", id, "UPDATE", userID, map[string]interface{}{"name": oldP.Name}, map[string]interface{}{"name": p.Name})
	return &p, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "project", id, "DELETE", userID, map[string]interface{}{"name": p.Name}, nil)
	return nil
}
