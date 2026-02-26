package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type workspaceService struct {
	repo     WorkspaceRepository
	auditLog AuditLogService
}

func NewWorkspaceService(repo WorkspaceRepository, auditLog AuditLogService) WorkspaceService {
	return &workspaceService{
		repo:     repo,
		auditLog: auditLog,
	}
}

func (s *workspaceService) GetWorkspace(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	return s.repo.Get(ctx, id)
}

func (s *workspaceService) ListWorkspaces(ctx context.Context) ([]domain.Workspace, error) {
	return s.repo.List(ctx)
}

func (s *workspaceService) CreateWorkspace(ctx context.Context, name, description string, userID uuid.UUID) (*domain.Workspace, error) {
	ws := &domain.Workspace{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
	}
	if err := s.repo.Create(ctx, ws); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "workspace", ws.ID, "CREATE", userID, nil, map[string]interface{}{"name": ws.Name})
	return ws, nil
}

func (s *workspaceService) UpdateWorkspace(ctx context.Context, id uuid.UUID, name, description string, userID uuid.UUID) (*domain.Workspace, error) {
	oldWS, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	ws := *oldWS
	ws.Name = name
	ws.Description = description
	if err := s.repo.Update(ctx, &ws); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "workspace", id, "UPDATE", userID, map[string]interface{}{"name": oldWS.Name}, map[string]interface{}{"name": ws.Name})
	return &ws, nil
}

func (s *workspaceService) DeleteWorkspace(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ws, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "workspace", id, "DELETE", userID, map[string]interface{}{"name": ws.Name}, nil)
	return nil
}
