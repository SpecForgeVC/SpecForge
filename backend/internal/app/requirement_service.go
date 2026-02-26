package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type requirementService struct {
	repo     RequirementRepository
	auditLog AuditLogService
}

func NewRequirementService(repo RequirementRepository, al AuditLogService) RequirementService {
	return &requirementService{repo: repo, auditLog: al}
}

func (s *requirementService) GetRequirement(ctx context.Context, id uuid.UUID) (*domain.Requirement, error) {
	return s.repo.Get(ctx, id)
}

func (s *requirementService) ListRequirements(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.Requirement, error) {
	return s.repo.List(ctx, roadmapItemID)
}

func (s *requirementService) CreateRequirement(ctx context.Context, roadmapItemID uuid.UUID, title, description string, testable bool, acceptanceCriteria string, orderIndex int, userID uuid.UUID) (*domain.Requirement, error) {
	req := &domain.Requirement{
		RoadmapItemID:      roadmapItemID,
		Title:              title,
		Description:        description,
		Testable:           testable,
		AcceptanceCriteria: acceptanceCriteria,
		OrderIndex:         orderIndex,
	}
	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "requirement", req.ID, "CREATE", userID, nil, map[string]interface{}{"title": title})
	return req, nil
}

func (s *requirementService) UpdateRequirement(ctx context.Context, id uuid.UUID, title, description string, testable bool, acceptanceCriteria string, orderIndex int, userID uuid.UUID) (*domain.Requirement, error) {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	req := &domain.Requirement{
		ID:                 id,
		RoadmapItemID:      old.RoadmapItemID,
		Title:              title,
		Description:        description,
		Testable:           testable,
		AcceptanceCriteria: acceptanceCriteria,
		OrderIndex:         orderIndex,
	}
	if err := s.repo.Update(ctx, req); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "requirement", id, "UPDATE", userID, map[string]interface{}{"title": old.Title}, map[string]interface{}{"title": title})
	return req, nil
}

func (s *requirementService) DeleteRequirement(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "requirement", id, "DELETE", userID, map[string]interface{}{"title": old.Title}, nil)
	return nil
}
