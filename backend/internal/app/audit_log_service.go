package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type auditLogService struct {
	repo AuditLogRepository
}

func NewAuditLogService(repo AuditLogRepository) AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) Log(ctx context.Context, entityType string, entityID uuid.UUID, action string, userID uuid.UUID, oldData, newData map[string]interface{}) error {
	log := &domain.AuditLog{
		ID:          uuid.New(),
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		PerformedBy: userID,
		OldData:     oldData,
		NewData:     newData,
	}
	return s.repo.Create(ctx, log)
}

func (s *auditLogService) GetEntityLogs(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error) {
	return s.repo.ListByEntity(ctx, entityType, entityID)
}

func (s *auditLogService) ListDriftEvents(ctx context.Context) ([]domain.AuditLog, error) {
	return s.repo.ListDriftEvents(ctx)
}
