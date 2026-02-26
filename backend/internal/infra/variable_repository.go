package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type variableRepository struct {
	queries *db.Queries
}

func NewVariableRepository(queries *db.Queries) app.VariableRepository {
	return &variableRepository{queries: queries}
}

func (r *variableRepository) Get(ctx context.Context, id uuid.UUID) (*domain.VariableDefinition, error) {
	row, err := r.queries.GetVariable(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapRow(row), nil
}

func (r *variableRepository) List(ctx context.Context, contractID uuid.UUID) ([]domain.VariableDefinition, error) {
	rows, err := r.queries.ListVariablesByContract(ctx, contractID)
	if err != nil {
		return nil, err
	}
	vars := make([]domain.VariableDefinition, len(rows))
	for i, row := range rows {
		vars[i] = *r.mapRow(row)
	}
	return vars, nil
}

func (r *variableRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VariableDefinition, error) {
	rows, err := r.queries.ListVariablesByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	vars := make([]domain.VariableDefinition, len(rows))
	for i, row := range rows {
		vars[i] = *r.mapRow(row)
	}
	return vars, nil
}

func (r *variableRepository) Create(ctx context.Context, v *domain.VariableDefinition) error {
	row, err := r.queries.CreateVariable(ctx, db.CreateVariableParams{
		ContractID:      v.ContractID,
		Name:            v.Name,
		Type:            v.Type,
		Required:        db.ToSqlBool(v.Required),
		DefaultValue:    db.ToSqlString(v.DefaultValue),
		Description:     db.ToSqlString(v.Description),
		ValidationRules: db.JSONToSql(v.ValidationRules),
	})
	if err != nil {
		return err
	}
	v.ID = row.ID
	return nil
}

func (r *variableRepository) Update(ctx context.Context, v *domain.VariableDefinition) error {
	_, err := r.queries.UpdateVariable(ctx, db.UpdateVariableParams{
		ID:              v.ID,
		Name:            v.Name,
		Type:            v.Type,
		Required:        db.ToSqlBool(v.Required),
		DefaultValue:    db.ToSqlString(v.DefaultValue),
		Description:     db.ToSqlString(v.Description),
		ValidationRules: db.JSONToSql(v.ValidationRules),
	})
	return err
}

func (r *variableRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteVariable(ctx, id)
}

func (r *variableRepository) mapRow(row db.VariableDefinition) *domain.VariableDefinition {
	return &domain.VariableDefinition{
		ID:              row.ID,
		ContractID:      row.ContractID,
		Name:            row.Name,
		Type:            row.Type,
		Required:        row.Required.Bool,
		DefaultValue:    row.DefaultValue.String,
		Description:     row.Description.String,
		ValidationRules: db.SqlToJSON(row.ValidationRules),
	}
}
