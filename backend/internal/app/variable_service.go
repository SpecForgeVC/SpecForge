package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type variableService struct {
	repo                VariableRepository
	contractRepo        ContractRepository
	roadmapRepo         RoadmapItemRepository
	auditLog            AuditLogService
	featureIntelligence FeatureIntelligenceService
	alignment           AlignmentService
}

func NewVariableService(repo VariableRepository, contractRepo ContractRepository, roadmapRepo RoadmapItemRepository, al AuditLogService, fi FeatureIntelligenceService, alignment AlignmentService) VariableService {
	return &variableService{
		repo:                repo,
		contractRepo:        contractRepo,
		roadmapRepo:         roadmapRepo,
		auditLog:            al,
		featureIntelligence: fi,
		alignment:           alignment,
	}
}

func (s *variableService) GetVariable(ctx context.Context, id uuid.UUID) (*domain.VariableDefinition, error) {
	return s.repo.Get(ctx, id)
}

func (s *variableService) ListVariables(ctx context.Context, contractID uuid.UUID) ([]domain.VariableDefinition, error) {
	return s.repo.List(ctx, contractID)
}

func (s *variableService) ListVariablesByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VariableDefinition, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *variableService) CreateVariable(ctx context.Context, contractID uuid.UUID, name, vType string, required bool, defaultValue, description string, validationRules map[string]interface{}, userID uuid.UUID) (*domain.VariableDefinition, error) {
	v := &domain.VariableDefinition{
		ContractID:      contractID,
		Name:            name,
		Type:            vType,
		Required:        required,
		DefaultValue:    defaultValue,
		Description:     description,
		ValidationRules: validationRules,
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "variable", v.ID, "CREATE", userID, nil, map[string]interface{}{"name": name})

	// Trigger intelligence recalculation
	contract, err := s.contractRepo.Get(ctx, contractID)
	if err == nil {
		_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, contract.RoadmapItemID)

		// Trigger alignment check
		if roadmapItem, err := s.roadmapRepo.Get(ctx, contract.RoadmapItemID); err == nil {
			_, _ = s.alignment.TriggerAlignmentCheck(ctx, roadmapItem.ProjectID)
		}
	}

	return v, nil
}

func (s *variableService) UpdateVariable(ctx context.Context, id uuid.UUID, name, vType string, required bool, defaultValue, description string, validationRules map[string]interface{}, userID uuid.UUID) (*domain.VariableDefinition, error) {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	v := &domain.VariableDefinition{
		ID:              id,
		ContractID:      old.ContractID,
		Name:            name,
		Type:            vType,
		Required:        required,
		DefaultValue:    defaultValue,
		Description:     description,
		ValidationRules: validationRules,
	}
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "variable", id, "UPDATE", userID, map[string]interface{}{"name": old.Name}, map[string]interface{}{"name": name})

	// Trigger intelligence recalculation
	contract, err := s.contractRepo.Get(ctx, old.ContractID)
	if err == nil {
		_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, contract.RoadmapItemID)

		// Trigger alignment check
		if roadmapItem, err := s.roadmapRepo.Get(ctx, contract.RoadmapItemID); err == nil {
			_, _ = s.alignment.TriggerAlignmentCheck(ctx, roadmapItem.ProjectID)
		}
	}

	return v, nil
}

func (s *variableService) DeleteVariable(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "variable", id, "DELETE", userID, map[string]interface{}{"name": old.Name}, nil)

	// Trigger intelligence recalculation
	contract, err := s.contractRepo.Get(ctx, old.ContractID)
	if err == nil {
		_, _ = s.featureIntelligence.CalculateFeatureScore(ctx, contract.RoadmapItemID)

		// Trigger alignment check
		if roadmapItem, err := s.roadmapRepo.Get(ctx, contract.RoadmapItemID); err == nil {
			_, _ = s.alignment.TriggerAlignmentCheck(ctx, roadmapItem.ProjectID)
		}
	}

	return nil
}
