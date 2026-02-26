package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type llmRepository struct {
	db db.DBTX
}

func NewLLMRepository(db db.DBTX) app.LLMRepository {
	return &llmRepository{db: db}
}

func (r *llmRepository) GetActive(ctx context.Context) (*domain.LLMConfiguration, error) {
	query := `
		SELECT id, provider, api_key, base_url, model, is_active, created_at, updated_at
		FROM llm_configurations
		ORDER BY is_active DESC, updated_at DESC
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query)

	var c domain.LLMConfiguration
	var baseURL sql.NullString
	var provider string

	err := row.Scan(
		&c.ID,
		&provider,
		&c.APIKey,
		&baseURL,
		&c.Model,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan llm config: %w", err)
	}
	c.Provider = domain.LLMProvider(provider)
	c.BaseURL = baseURL.String

	return &c, nil
}

func (r *llmRepository) Upsert(ctx context.Context, config *domain.LLMConfiguration) error {
	// If setting active, unset others (handled by unique index roughly, but cleaner to update)
	if config.IsActive {
		if _, err := r.db.ExecContext(ctx, "UPDATE llm_configurations SET is_active = false"); err != nil {
			return err
		}
	}

	if config.ID == uuid.Nil {
		config.ID = uuid.New()
	}

	query := `
		INSERT INTO llm_configurations (id, provider, api_key, base_url, model, is_active, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (id) DO UPDATE SET
			provider = EXCLUDED.provider,
			api_key = EXCLUDED.api_key,
			base_url = EXCLUDED.base_url,
			model = EXCLUDED.model,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		config.ID,
		config.Provider,
		config.APIKey,
		sql.NullString{String: config.BaseURL, Valid: config.BaseURL != ""},
		config.Model,
		config.IsActive,
	)
	return err
}
