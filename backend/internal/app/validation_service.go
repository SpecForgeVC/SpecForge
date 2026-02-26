package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type validationRuleService struct {
	repo     ValidationRuleRepository
	auditLog AuditLogService
}

func NewValidationRuleService(repo ValidationRuleRepository, al AuditLogService) ValidationRuleService {
	return &validationRuleService{repo: repo, auditLog: al}
}

func (s *validationRuleService) GetValidationRule(ctx context.Context, id uuid.UUID) (*domain.ValidationRule, error) {
	return s.repo.Get(ctx, id)
}

func (s *validationRuleService) ListValidationRules(ctx context.Context, projectID uuid.UUID) ([]domain.ValidationRule, error) {
	return s.repo.List(ctx, projectID)
}

func (s *validationRuleService) CreateValidationRule(ctx context.Context, projectID uuid.UUID, name, rType string, config map[string]interface{}, description string, userID uuid.UUID) (*domain.ValidationRule, error) {
	rule := &domain.ValidationRule{
		ProjectID:   projectID,
		Name:        name,
		RuleType:    rType,
		RuleConfig:  config,
		Description: description,
	}
	if err := s.repo.Create(ctx, rule); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "validation_rule", rule.ID, "CREATE", userID, nil, map[string]interface{}{"name": name})
	return rule, nil
}

func (s *validationRuleService) UpdateValidationRule(ctx context.Context, id uuid.UUID, name, rType string, config map[string]interface{}, description string, userID uuid.UUID) (*domain.ValidationRule, error) {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	rule := &domain.ValidationRule{
		ID:          id,
		ProjectID:   old.ProjectID,
		Name:        name,
		RuleType:    rType,
		RuleConfig:  config,
		Description: description,
	}
	if err := s.repo.Update(ctx, rule); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "validation_rule", id, "UPDATE", userID, map[string]interface{}{"name": old.Name}, map[string]interface{}{"name": name})
	return rule, nil
}

func (s *validationRuleService) DeleteValidationRule(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "validation_rule", id, "DELETE", userID, map[string]interface{}{"name": old.Name}, nil)
	return nil
}
