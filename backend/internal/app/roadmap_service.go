package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type roadmapItemService struct {
	repo                RoadmapItemRepository
	auditLog            AuditLogService
	featureIntelligence FeatureIntelligenceService
	governance          GovernanceService
	alignment           AlignmentService
}

func NewRoadmapItemService(repo RoadmapDependencyRepository, roadmapRepo RoadmapItemRepository, auditLog AuditLogService, fi FeatureIntelligenceService, gov GovernanceService, alignment AlignmentService) RoadmapItemService {
	return &roadmapItemService{
		repo:                roadmapRepo,
		auditLog:            auditLog,
		featureIntelligence: fi,
		governance:          gov,
		alignment:           alignment,
	}
}

func (s *roadmapItemService) GetRoadmapItem(ctx context.Context, id uuid.UUID) (*domain.RoadmapItem, error) {
	return s.repo.Get(ctx, id)
}

func (s *roadmapItemService) ListRoadmapItems(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapItem, error) {
	return s.repo.List(ctx, projectID)
}

func (s *roadmapItemService) CreateRoadmapItem(ctx context.Context, item *domain.RoadmapItem, userID uuid.UUID) (*domain.RoadmapItem, error) {
	item.ID = uuid.New()
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "roadmap_item", item.ID, "CREATE", userID, nil, map[string]interface{}{"title": item.Title})

	// Trigger intelligence recalculation
	if item.Type == domain.Feature {
		_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, item.ID)
	}

	// Trigger alignment check
	_, _ = s.alignment.TriggerAlignmentCheck(ctx, item.ProjectID)

	return item, nil
}

func (s *roadmapItemService) UpdateRoadmapItem(ctx context.Context, id uuid.UUID, title, description, businessContext, technicalContext string, status domain.RoadmapItemStatus, userID uuid.UUID) (*domain.RoadmapItem, error) {
	oldItem, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	item := *oldItem
	item.Title = title
	item.Description = description
	item.BusinessContext = businessContext
	item.TechnicalContext = technicalContext

	// Governance Check on Status Transition
	if status == domain.StatusComplete || status == domain.StatusInProgress {
		canBuild, reasons, err := s.governance.CanBuildFeature(ctx, id)
		if err != nil {
			return nil, err
		}
		if !canBuild {
			return nil, fmt.Errorf("governance check failed: %v", reasons)
		}
	}

	item.Status = status
	if err := s.repo.Update(ctx, &item); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "roadmap_item", id, "UPDATE", userID, map[string]interface{}{"status": oldItem.Status}, map[string]interface{}{"status": item.Status})

	// Trigger intelligence recalculation
	if item.Type == domain.Feature {
		_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, id)
	}

	// Trigger alignment check
	_, _ = s.alignment.TriggerAlignmentCheck(ctx, item.ProjectID)

	return &item, nil
}

func (s *roadmapItemService) DeleteRoadmapItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "roadmap_item", id, "DELETE", userID, map[string]interface{}{"title": item.Title}, nil)
	return nil
}
