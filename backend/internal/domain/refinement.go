package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefinementStatus string

const (
	RefinementStatusInProgress RefinementStatus = "IN_PROGRESS"
	RefinementStatusValidated  RefinementStatus = "VALIDATED"
	RefinementStatusFailed     RefinementStatus = "FAILED"
	RefinementStatusApproved   RefinementStatus = "APPROVED"
)

type RefinementSession struct {
	ID               uuid.UUID        `json:"id"`
	ArtifactType     string           `json:"artifact_type"`
	TargetType       string           `json:"target_type"`
	InitialPrompt    string           `json:"initial_prompt"`
	ContextData      map[string]any   `json:"context_data"`
	MaxIterations    int              `json:"max_iterations"`
	CurrentIteration int              `json:"current_iteration"`
	Status           RefinementStatus `json:"status"`
	ConfidenceScore  float64          `json:"confidence_score"`
	ValidationErrors []string         `json:"validation_errors"`
	Result           map[string]any   `json:"result"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type RefinementIteration struct {
	ID               uuid.UUID             `json:"id"`
	SessionID        uuid.UUID             `json:"session_id"`
	Iteration        int                   `json:"iteration"`
	Prompt           string                `json:"prompt"`
	Response         string                `json:"response"`
	Artifact         map[string]any        `json:"artifact"`
	ValidationResult map[string]any        `json:"validation_result"`
	SelfEvaluation   *SelfEvaluationResult `json:"self_evaluation,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
}

type SelfEvaluationResult struct {
	Score                  int      `json:"score"`
	AmbiguityFlags         []string `json:"ambiguity_flags"`
	MissingConstraints     []string `json:"missing_constraints"`
	WeakValidations        []string `json:"weak_validations"`
	SecurityConcerns       []string `json:"security_concerns"`
	ImprovementSuggestions []string `json:"improvement_suggestions"`
}

type RefinementEvent struct {
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Payload map[string]any `json:"payload,omitempty"`
}
