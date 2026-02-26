package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type MCPTokenService interface {
	GenerateToken(ctx context.Context, userID, projectID uuid.UUID) (*domain.MCPTokenCreate, error)
	ValidateToken(ctx context.Context, rawToken string) (*domain.MCPToken, error)
	ListTokens(ctx context.Context, projectID uuid.UUID) ([]domain.MCPToken, error)
	RevokeToken(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllForProject(ctx context.Context, projectID uuid.UUID) error
}

type mcpTokenService struct {
	repo MCPTokenRepository
}

func NewMCPTokenService(repo MCPTokenRepository) MCPTokenService {
	return &mcpTokenService{repo: repo}
}

func (s *mcpTokenService) GenerateToken(ctx context.Context, userID, projectID uuid.UUID) (*domain.MCPTokenCreate, error) {
	// Generate 64-character secure random token
	b := make([]byte, 32) // 32 bytes = 64 hex chars
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random token: %w", err)
	}
	rawToken := "sf_live_" + hex.EncodeToString(b)
	prefix := rawToken[:11] // "sf_live_" + first 3 chars of hex

	// Hash the token
	hash := s.hashToken(rawToken)

	token := &domain.MCPToken{
		ID:          uuid.New(),
		UserID:      userID,
		ProjectID:   projectID,
		TokenHash:   hash,
		TokenPrefix: prefix,
		Revoked:     false,
	}

	if err := s.repo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to store token handle: %w", err)
	}

	return &domain.MCPTokenCreate{
		UserID:    userID,
		ProjectID: projectID,
		TokenRaw:  rawToken,
	}, nil
}

func (s *mcpTokenService) ValidateToken(ctx context.Context, rawToken string) (*domain.MCPToken, error) {
	hash := s.hashToken(rawToken)
	token, err := s.repo.GetByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, fmt.Errorf("invalid token")
	}

	// Update last used time asynchronously or synchronously
	_ = s.repo.UpdateUsage(ctx, token.ID)

	return token, nil
}

func (s *mcpTokenService) ListTokens(ctx context.Context, projectID uuid.UUID) ([]domain.MCPToken, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *mcpTokenService) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	return s.repo.Revoke(ctx, tokenID)
}

func (s *mcpTokenService) RevokeAllForProject(ctx context.Context, projectID uuid.UUID) error {
	return s.repo.RevokeAllForProject(ctx, projectID)
}

func (s *mcpTokenService) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	// In a real production app, we should use a salt, which could be the projectID or a dedicated column.
	// For simplicity and matching common API token hashing patterns, we'll use a straight SHA-256 here
	// but strictly speaking, Argon2 or Bcrypt is better for passwords.
	// API tokens are often high-entropy enough that SHA-256 is acceptable if the tokens are long.
	return hex.EncodeToString(h.Sum(nil))
}
