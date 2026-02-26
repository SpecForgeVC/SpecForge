package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type auditLogRepository struct {
	queries *db.Queries
}

func NewAuditLogRepository(queries *db.Queries) app.AuditLogRepository {
	return &auditLogRepository{queries: queries}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	_, err := r.queries.CreateAuditLog(ctx, db.CreateAuditLogParams{
		EntityType:  log.EntityType,
		EntityID:    log.EntityID,
		Action:      log.Action,
		PerformedBy: uuid.NullUUID{UUID: log.PerformedBy, Valid: log.PerformedBy != uuid.Nil},
		OldData:     db.JSONToSql(log.OldData),
		NewData:     db.JSONToSql(log.NewData),
	})
	return err
}

func (r *auditLogRepository) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	rows, err := r.queries.ListAuditLogsByEntity(ctx, db.ListAuditLogsByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, err
	}
	return r.mapRows(rows), nil
}

func (r *auditLogRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.AuditLog, error) {
	rows, err := r.queries.ListAuditLogsByUser(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	return r.mapRows(rows), nil
}

func (r *auditLogRepository) ListDriftEvents(ctx context.Context) ([]domain.AuditLog, error) {
	rows, err := r.queries.ListAuditLogsByAction(ctx, "DRIFT_DETECTED")
	if err != nil {
		return nil, err
	}
	return r.mapRows(rows), nil
}

func (r *auditLogRepository) mapRows(rows []db.AuditLog) []domain.AuditLog {
	logs := make([]domain.AuditLog, len(rows))
	for i, row := range rows {
		var userID uuid.UUID
		if row.PerformedBy.Valid {
			userID = row.PerformedBy.UUID
		}
		logs[i] = domain.AuditLog{
			ID:          row.ID,
			EntityType:  row.EntityType,
			EntityID:    row.EntityID,
			Action:      row.Action,
			PerformedBy: userID,
			OldData:     db.SqlToJSON(row.OldData),
			NewData:     db.SqlToJSON(row.NewData),
			CreatedAt:   row.CreatedAt.Time,
		}
	}
	return logs
}
