package infra

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type MCPTokenRepository struct {
	db *sql.DB
}

func NewMCPTokenRepository(db *sql.DB) app.MCPTokenRepository {
	return &MCPTokenRepository{db: db}
}

func (r *MCPTokenRepository) Create(ctx context.Context, token *domain.MCPToken) error {
	query := `
		INSERT INTO mcp_api_tokens (id, user_id, project_id, token_hash, token_prefix, revoked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	token.CreatedAt = now
	token.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.ProjectID, token.TokenHash, token.TokenPrefix, token.Revoked, token.CreatedAt, token.UpdatedAt,
	)
	return err
}

func (r *MCPTokenRepository) GetByHash(ctx context.Context, hash string) (*domain.MCPToken, error) {
	query := `
		SELECT id, user_id, project_id, token_hash, token_prefix, revoked, last_used_at, created_at, updated_at
		FROM mcp_api_tokens
		WHERE token_hash = $1 AND revoked = FALSE
	`
	var t domain.MCPToken
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&t.ID, &t.UserID, &t.ProjectID, &t.TokenHash, &t.TokenPrefix, &t.Revoked, &t.LastUsedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *MCPTokenRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.MCPToken, error) {
	query := `
		SELECT id, user_id, project_id, token_hash, token_prefix, revoked, last_used_at, created_at, updated_at
		FROM mcp_api_tokens
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.MCPToken
	for rows.Next() {
		var t domain.MCPToken
		err := rows.Scan(
			&t.ID, &t.UserID, &t.ProjectID, &t.TokenHash, &t.TokenPrefix, &t.Revoked, &t.LastUsedAt, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (r *MCPTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE mcp_api_tokens SET revoked = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *MCPTokenRepository) RevokeAllForProject(ctx context.Context, projectID uuid.UUID) error {
	query := `UPDATE mcp_api_tokens SET revoked = TRUE, updated_at = CURRENT_TIMESTAMP WHERE project_id = $1`
	_, err := r.db.ExecContext(ctx, query, projectID)
	return err
}

func (r *MCPTokenRepository) UpdateUsage(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE mcp_api_tokens SET last_used_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *MCPTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM mcp_api_tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
