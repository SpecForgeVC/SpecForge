package infra

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
	"github.com/sqlc-dev/pqtype"
)

type variableLineageRepository struct {
	queries *db.Queries
}

func NewVariableLineageRepository(queries *db.Queries) app.VariableLineageRepository {
	return &variableLineageRepository{queries: queries}
}

func (r *variableLineageRepository) CreateEvent(ctx context.Context, event *domain.VariableLineageEvent) error {
	metaBytes, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	_, err = r.queries.CreateVariableLineageEvent(ctx, db.CreateVariableLineageEventParams{
		VariableID:      event.VariableID,
		EventType:       string(event.EventType),
		SourceComponent: event.SourceComponent,
		Description:     event.Description,
		PerformedBy:     uuid.NullUUID{UUID: event.PerformedBy, Valid: event.PerformedBy != uuid.Nil},
		Metadata:        pqtype.NullRawMessage{RawMessage: metaBytes, Valid: len(metaBytes) > 0},
	})
	return err
}

func (r *variableLineageRepository) ListEvents(ctx context.Context, variableID uuid.UUID) ([]domain.VariableLineageEvent, error) {
	rows, err := r.queries.GetVariableLineageEvents(ctx, variableID)
	if err != nil {
		return nil, err
	}

	events := make([]domain.VariableLineageEvent, len(rows))
	for i, row := range rows {
		var meta map[string]interface{}
		if row.Metadata.Valid {
			_ = json.Unmarshal(row.Metadata.RawMessage, &meta)
		}

		events[i] = domain.VariableLineageEvent{
			ID:              row.ID,
			VariableID:      row.VariableID,
			EventType:       domain.LineageEventType(fmt.Sprintf("%s", row.EventType)),
			SourceComponent: row.SourceComponent,
			Description:     row.Description,
			PerformedBy:     row.PerformedBy.UUID,
			CreatedAt:       row.CreatedAt.Time,
			Metadata:        meta,
		}
	}
	return events, nil
}

func (r *variableLineageRepository) CreateDependency(ctx context.Context, dep *domain.VariableDependency) error {
	_, err := r.queries.CreateVariableDependency(ctx, db.CreateVariableDependencyParams{
		SourceVariableID: dep.SourceVariableID,
		TargetVariableID: dep.TargetVariableID,
		DependencyType:   string(dep.DependencyType),
	})
	return err
}

func (r *variableLineageRepository) ListDependencies(ctx context.Context, variableID uuid.UUID) ([]domain.VariableDependency, error) {
	rows, err := r.queries.GetVariableDependencies(ctx, variableID)
	if err != nil {
		return nil, err
	}

	deps := make([]domain.VariableDependency, len(rows))
	for i, row := range rows {
		deps[i] = domain.VariableDependency{
			ID:               row.ID,
			SourceVariableID: row.SourceVariableID,
			TargetVariableID: row.TargetVariableID,
			DependencyType:   domain.DependencyType(fmt.Sprintf("%s", row.DependencyType)),
			CreatedAt:        row.CreatedAt.Time,
		}
	}
	return deps, nil
}
