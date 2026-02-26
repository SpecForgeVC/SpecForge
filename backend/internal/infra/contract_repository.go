package infra

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type contractRepository struct {
	queries *db.Queries
}

func NewContractRepository(queries *db.Queries) app.ContractRepository {
	return &contractRepository{queries: queries}
}

func (r *contractRepository) Get(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error) {
	row, err := r.queries.GetContractDefinition(ctx, id)
	if err != nil {
		return nil, err
	}
	var deprecatedFields []string
	if row.DeprecatedFields.Valid {
		if err := json.Unmarshal(row.DeprecatedFields.RawMessage, &deprecatedFields); err != nil {
			deprecatedFields = []string{}
		}
	} else {
		deprecatedFields = []string{}
	}

	return &domain.ContractDefinition{
		ID:                 row.ID,
		RoadmapItemID:      row.RoadmapItemID,
		ContractType:       domain.ContractType(row.ContractType),
		Version:            row.Version,
		InputSchema:        db.SqlToJSON(row.InputSchema),
		OutputSchema:       db.SqlToJSON(row.OutputSchema),
		ErrorSchema:        db.SqlToJSON(row.ErrorSchema),
		BackwardCompatible: row.BackwardCompatible.Bool,
		DeprecatedFields:   deprecatedFields,
		CreatedAt:          row.CreatedAt.Time,
	}, nil
}

func (r *contractRepository) List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.ContractDefinition, error) {
	rows, err := r.queries.ListContractDefinitions(ctx, roadmapItemID)
	if err != nil {
		return nil, err
	}
	contracts := make([]domain.ContractDefinition, len(rows))
	for i, row := range rows {
		contracts[i] = *r.mapRow(row)
	}
	return contracts, nil
}

func (r *contractRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.ContractDefinition, error) {
	rows, err := r.queries.ListContractDefinitionsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	contracts := make([]domain.ContractDefinition, len(rows))
	for i, row := range rows {
		contracts[i] = *r.mapRow(row)
	}
	return contracts, nil
}

func (r *contractRepository) mapRow(row db.ContractDefinition) *domain.ContractDefinition {
	var deprecatedFields []string
	if row.DeprecatedFields.Valid {
		if err := json.Unmarshal(row.DeprecatedFields.RawMessage, &deprecatedFields); err != nil {
			deprecatedFields = []string{}
		}
	} else {
		deprecatedFields = []string{}
	}

	return &domain.ContractDefinition{
		ID:                 row.ID,
		RoadmapItemID:      row.RoadmapItemID,
		ContractType:       domain.ContractType(row.ContractType),
		Version:            row.Version,
		InputSchema:        db.SqlToJSON(row.InputSchema),
		OutputSchema:       db.SqlToJSON(row.OutputSchema),
		ErrorSchema:        db.SqlToJSON(row.ErrorSchema),
		BackwardCompatible: row.BackwardCompatible.Bool,
		DeprecatedFields:   deprecatedFields,
		CreatedAt:          row.CreatedAt.Time,
	}
}

func (r *contractRepository) Create(ctx context.Context, c *domain.ContractDefinition) error {
	deprecatedFields, _ := json.Marshal(c.DeprecatedFields)
	row, err := r.queries.CreateContractDefinition(ctx, db.CreateContractDefinitionParams{
		RoadmapItemID:      c.RoadmapItemID,
		ContractType:       db.ContractType(c.ContractType),
		Version:            c.Version,
		InputSchema:        db.JSONToSql(c.InputSchema),
		OutputSchema:       db.JSONToSql(c.OutputSchema),
		ErrorSchema:        db.JSONToSql(c.ErrorSchema),
		BackwardCompatible: db.BoolToSql(c.BackwardCompatible),
		DeprecatedFields:   db.BytesToPQRawMessage(deprecatedFields),
	})
	if err != nil {
		return err
	}

	c.ID = row.ID
	c.CreatedAt = row.CreatedAt.Time
	return nil
}

func (r *contractRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteContractDefinition(ctx, id)
}

func (r *contractRepository) Update(ctx context.Context, c *domain.ContractDefinition) error {
	deprecatedFields, _ := json.Marshal(c.DeprecatedFields)
	_, err := r.queries.UpdateContractDefinition(ctx, db.UpdateContractDefinitionParams{
		ID:                 c.ID,
		ContractType:       db.ContractType(c.ContractType),
		Version:            c.Version,
		InputSchema:        db.JSONToSql(c.InputSchema),
		OutputSchema:       db.JSONToSql(c.OutputSchema),
		ErrorSchema:        db.JSONToSql(c.ErrorSchema),
		BackwardCompatible: db.BoolToSql(c.BackwardCompatible),
		DeprecatedFields:   db.BytesToPQRawMessage(deprecatedFields),
	})
	return err
}
