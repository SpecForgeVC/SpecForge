package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type validationRuleRepository struct {
	queries *db.Queries
}

func NewValidationRuleRepository(queries *db.Queries) app.ValidationRuleRepository {
	return &validationRuleRepository{queries: queries}
}

func (r *validationRuleRepository) Get(ctx context.Context, id uuid.UUID) (*domain.ValidationRule, error) {
	row, err := r.queries.GetValidationRule(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapRow(row), nil
}

func (r *validationRuleRepository) List(ctx context.Context, projectID uuid.UUID) ([]domain.ValidationRule, error) {
	rows, err := r.queries.ListValidationRulesByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	rules := make([]domain.ValidationRule, len(rows))
	for i, row := range rows {
		rules[i] = *r.mapRow(row)
	}
	return rules, nil
}

func (r *validationRuleRepository) Create(ctx context.Context, rule *domain.ValidationRule) error {
	row, err := r.queries.CreateValidationRule(ctx, db.CreateValidationRuleParams{
		ProjectID:   rule.ProjectID,
		Name:        rule.Name,
		RuleType:    rule.RuleType,
		RuleConfig:  db.JSONToSql(rule.RuleConfig),
		Description: db.ToSqlString(rule.Description),
	})
	if err != nil {
		return err
	}
	rule.ID = row.ID
	return nil
}

func (r *validationRuleRepository) Update(ctx context.Context, rule *domain.ValidationRule) error {
	_, err := r.queries.UpdateValidationRule(ctx, db.UpdateValidationRuleParams{
		ID:          rule.ID,
		Name:        rule.Name,
		RuleType:    rule.RuleType,
		RuleConfig:  db.JSONToSql(rule.RuleConfig),
		Description: db.ToSqlString(rule.Description),
	})
	return err
}

func (r *validationRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteValidationRule(ctx, id)
}

func (r *validationRuleRepository) mapRow(row db.ValidationRule) *domain.ValidationRule {
	return &domain.ValidationRule{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		RuleType:    row.RuleType,
		RuleConfig:  db.SqlToJSON(row.RuleConfig),
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt.Time,
	}
}
