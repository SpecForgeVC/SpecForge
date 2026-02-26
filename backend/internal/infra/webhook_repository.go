package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type webhookRepository struct {
	queries *db.Queries
}

func NewWebhookRepository(queries *db.Queries) app.WebhookRepository {
	return &webhookRepository{queries: queries}
}

func (r *webhookRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Webhook, error) {
	row, err := r.queries.GetWebhook(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.mapRow(row), nil
}

func (r *webhookRepository) List(ctx context.Context, projectID uuid.UUID) ([]domain.Webhook, error) {
	rows, err := r.queries.ListWebhooksByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	webhooks := make([]domain.Webhook, len(rows))
	for i, row := range rows {
		webhooks[i] = *r.mapRow(row)
	}
	return webhooks, nil
}

func (r *webhookRepository) Create(ctx context.Context, w *domain.Webhook) error {
	row, err := r.queries.CreateWebhook(ctx, db.CreateWebhookParams{
		ProjectID: w.ProjectID,
		Url:       w.URL,
		Events:    w.Events,
		Secret:    w.Secret,
		Active:    db.ToSqlBool(w.Active),
	})
	if err != nil {
		return err
	}
	w.ID = row.ID
	return nil
}

func (r *webhookRepository) Update(ctx context.Context, w *domain.Webhook) error {
	_, err := r.queries.UpdateWebhook(ctx, db.UpdateWebhookParams{
		ID:     w.ID,
		Url:    w.URL,
		Events: w.Events,
		Secret: w.Secret,
		Active: db.ToSqlBool(w.Active),
	})
	return err
}

func (r *webhookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteWebhook(ctx, id)
}

func (r *webhookRepository) mapRow(row db.Webhook) *domain.Webhook {
	return &domain.Webhook{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		URL:       row.Url,
		Events:    row.Events,
		Secret:    row.Secret,
		Active:    row.Active.Bool,
		CreatedAt: row.CreatedAt.Time,
	}
}
