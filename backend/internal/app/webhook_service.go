package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type webhookService struct {
	repo     WebhookRepository
	auditLog AuditLogService
}

func NewWebhookService(repo WebhookRepository, al AuditLogService) WebhookService {
	return &webhookService{repo: repo, auditLog: al}
}

func (s *webhookService) GetWebhook(ctx context.Context, id uuid.UUID) (*domain.Webhook, error) {
	return s.repo.Get(ctx, id)
}

func (s *webhookService) ListWebhooks(ctx context.Context, projectID uuid.UUID) ([]domain.Webhook, error) {
	return s.repo.List(ctx, projectID)
}

func (s *webhookService) CreateWebhook(ctx context.Context, projectID uuid.UUID, url string, events []string, secret string, userID uuid.UUID) (*domain.Webhook, error) {
	w := &domain.Webhook{
		ProjectID: projectID,
		URL:       url,
		Events:    events,
		Secret:    secret,
		Active:    true,
	}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "webhook", w.ID, "CREATE", userID, nil, map[string]interface{}{"url": url})
	return w, nil
}

func (s *webhookService) UpdateWebhook(ctx context.Context, id uuid.UUID, url string, events []string, secret string, active bool, userID uuid.UUID) (*domain.Webhook, error) {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	w := &domain.Webhook{
		ID:        id,
		ProjectID: old.ProjectID,
		URL:       url,
		Events:    events,
		Secret:    secret,
		Active:    active,
	}
	if err := s.repo.Update(ctx, w); err != nil {
		return nil, err
	}
	s.auditLog.Log(ctx, "webhook", id, "UPDATE", userID, map[string]interface{}{"url": old.URL}, map[string]interface{}{"url": url})
	return w, nil
}

func (s *webhookService) DeleteWebhook(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	old, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "webhook", id, "DELETE", userID, map[string]interface{}{"url": old.URL}, nil)
	return nil
}
