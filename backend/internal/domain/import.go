package domain

import (
	"time"

	"github.com/google/uuid"
)

type ImportStatus string

const (
	ImportStatusPartial  ImportStatus = "partial"
	ImportStatusComplete ImportStatus = "complete"
)

type ImportSession struct {
	ID                uuid.UUID    `json:"id"`
	ProjectID         uuid.UUID    `json:"project_id"`
	CompletenessScore int32        `json:"completeness_score"`
	Status            ImportStatus `json:"status"`
	IterationCount    int32        `json:"iteration_count"`
	Locked            bool         `json:"locked"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type ImportArtifact struct {
	ID        uuid.UUID              `json:"id"`
	SessionID uuid.UUID              `json:"session_id"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}
