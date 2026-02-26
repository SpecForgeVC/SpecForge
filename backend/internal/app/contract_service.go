package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type contractService struct {
	repo                ContractRepository
	roadmapRepo         RoadmapItemRepository
	featureIntelligence FeatureIntelligenceService
	governance          GovernanceService
	alignment           AlignmentService
}

func NewContractService(repo ContractRepository, roadmapRepo RoadmapItemRepository, fi FeatureIntelligenceService, gov GovernanceService, alignment AlignmentService) ContractService {
	return &contractService{
		repo:                repo,
		roadmapRepo:         roadmapRepo,
		featureIntelligence: fi,
		governance:          gov,
		alignment:           alignment,
	}
}

func (s *contractService) GetContract(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error) {
	return s.repo.Get(ctx, id)
}

func (s *contractService) ListContracts(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.ContractDefinition, error) {
	return s.repo.List(ctx, roadmapItemID)
}

func (s *contractService) ListContractsByProject(ctx context.Context, projectID uuid.UUID) ([]domain.ContractDefinition, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *contractService) CreateContract(ctx context.Context, roadmapItemID uuid.UUID, cType domain.ContractType, version string, input, output, errSchema map[string]interface{}) (*domain.ContractDefinition, error) {
	c := &domain.ContractDefinition{
		ID:                 uuid.New(),
		RoadmapItemID:      roadmapItemID,
		ContractType:       cType,
		Version:            version,
		InputSchema:        input,
		OutputSchema:       output,
		ErrorSchema:        errSchema,
		BackwardCompatible: true,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	// Trigger intelligence recalculation
	_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, roadmapItemID)

	// Trigger alignment check
	if roadmapItem, err := s.roadmapRepo.Get(ctx, roadmapItemID); err == nil {
		_, _ = s.alignment.TriggerAlignmentCheck(ctx, roadmapItem.ProjectID)
	}

	return c, nil
}

func (s *contractService) DeleteContract(ctx context.Context, id uuid.UUID) error {
	// Need to get the contract first to know the feature ID
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	// Trigger intelligence recalculation
	_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, c.RoadmapItemID)

	// Trigger alignment check
	if roadmapItem, err := s.roadmapRepo.Get(ctx, c.RoadmapItemID); err == nil {
		_, _ = s.alignment.TriggerAlignmentCheck(ctx, roadmapItem.ProjectID)
	}

	return nil
}

func (s *contractService) UpdateContract(ctx context.Context, id uuid.UUID, cType domain.ContractType, version string, input, output, errSchema map[string]interface{}) (*domain.ContractDefinition, error) {
	// Governance Check
	allowed, reasons, err := s.governance.CanUpdateContract(ctx, id)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("governance check failed: %v", reasons)
	}

	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	c := &domain.ContractDefinition{
		ID:                 id,
		RoadmapItemID:      old.RoadmapItemID,
		ContractType:       cType,
		Version:            version,
		InputSchema:        input,
		OutputSchema:       output,
		ErrorSchema:        errSchema,
		BackwardCompatible: old.BackwardCompatible,
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}

	// Trigger intelligence recalculation
	_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, old.RoadmapItemID)

	// Trigger alignment check
	if roadmapItem, err := s.roadmapRepo.Get(ctx, old.RoadmapItemID); err == nil {
		_, _ = s.alignment.TriggerAlignmentCheck(ctx, roadmapItem.ProjectID)
	}

	return c, nil
}
