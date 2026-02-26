package domain

import (
	"time"

	"github.com/google/uuid"
)

type MCPToken struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	TokenHash   string     `json:"-"` // Never expose hash in JSON
	TokenPrefix string     `json:"token_prefix"`
	Revoked     bool       `json:"revoked"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type MCPTokenCreate struct {
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
	TokenRaw  string    `json:"token_raw"` // Used during generation to return to user
}

type MCPTokenFilter struct {
	UserID    *uuid.UUID
	ProjectID *uuid.UUID
	Revoked   *bool
}
