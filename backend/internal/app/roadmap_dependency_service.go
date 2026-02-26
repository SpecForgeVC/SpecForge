package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type roadmapDependencyService struct {
	repo        RoadmapDependencyRepository
	roadmapRepo RoadmapItemRepository
	auditLog    AuditLogService
}

func NewRoadmapDependencyService(repo RoadmapDependencyRepository, roadmapRepo RoadmapItemRepository, auditLog AuditLogService) RoadmapDependencyService {
	return &roadmapDependencyService{
		repo:        repo,
		roadmapRepo: roadmapRepo,
		auditLog:    auditLog,
	}
}

func (s *roadmapDependencyService) CreateDependency(ctx context.Context, sourceID, targetID uuid.UUID, dType domain.DependencyType) (*domain.RoadmapDependency, error) {
	// Verify both items exist
	if _, err := s.roadmapRepo.Get(ctx, sourceID); err != nil {
		return nil, err
	}
	if _, err := s.roadmapRepo.Get(ctx, targetID); err != nil {
		return nil, err
	}

	dep := &domain.RoadmapDependency{
		ID:             uuid.New(),
		SourceID:       sourceID,
		TargetID:       targetID,
		DependencyType: dType,
	}

	if err := s.repo.Create(ctx, dep); err != nil {
		return nil, err
	}

	// Log audit
	s.auditLog.Log(ctx, "roadmap_dependency", dep.ID, "CREATE", uuid.Nil, nil, map[string]interface{}{"source_id": sourceID, "target_id": targetID})

	return dep, nil
}

func (s *roadmapDependencyService) ListDependencies(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapDependency, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *roadmapDependencyService) DeleteDependency(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
