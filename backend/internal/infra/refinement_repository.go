package infra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type refinementRepository struct {
	db db.DBTX
}

func NewRefinementRepository(db db.DBTX) app.RefinementRepository {
	return &refinementRepository{db: db}
}

func (r *refinementRepository) CreateSession(ctx context.Context, s *domain.RefinementSession) error {
	query := `
		INSERT INTO refinement_sessions 
		(id, artifact_type, initial_prompt, context_data, max_iterations, current_iteration, status, confidence_score, validation_errors, result, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	contextJSON, _ := json.Marshal(s.ContextData)
	resultJSON, _ := json.Marshal(s.Result)
	if s.Result == nil {
		resultJSON = nil // Use NULL in DB
	}

	_, err := r.db.ExecContext(ctx, query,
		s.ID,
		s.ArtifactType,
		s.InitialPrompt,
		contextJSON,
		s.MaxIterations,
		s.CurrentIteration,
		s.Status,
		s.ConfidenceScore,
		pq.Array(s.ValidationErrors),
		resultJSON,
		s.CreatedAt,
		s.UpdatedAt,
	)
	return err
}

func (r *refinementRepository) GetSession(ctx context.Context, id uuid.UUID) (*domain.RefinementSession, error) {
	query := `
		SELECT id, artifact_type, initial_prompt, context_data, max_iterations, current_iteration, status, confidence_score, validation_errors, result, created_at, updated_at
		FROM refinement_sessions
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.RefinementSession
	var contextData, result []byte
	var validationErrors []string
	var status string
	var confidenceScore sql.NullFloat64

	err := row.Scan(
		&s.ID,
		&s.ArtifactType,
		&s.InitialPrompt,
		&contextData,
		&s.MaxIterations,
		&s.CurrentIteration,
		&status,
		&confidenceScore,
		pq.Array(&validationErrors),
		&result,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	s.Status = domain.RefinementStatus(status)
	s.ConfidenceScore = confidenceScore.Float64
	s.ValidationErrors = validationErrors

	if len(contextData) > 0 {
		_ = json.Unmarshal(contextData, &s.ContextData)
	}
	if len(result) > 0 {
		_ = json.Unmarshal(result, &s.Result)
	}

	return &s, nil
}

func (r *refinementRepository) UpdateSession(ctx context.Context, s *domain.RefinementSession) error {
	query := `
		UPDATE refinement_sessions
		SET current_iteration = $2, status = $3, confidence_score = $4, validation_errors = $5, result = $6, updated_at = NOW()
		WHERE id = $1
	`
	resultJSON, _ := json.Marshal(s.Result)
	if s.Result == nil {
		resultJSON = nil
	}

	_, err := r.db.ExecContext(ctx, query,
		s.ID,
		s.CurrentIteration,
		s.Status,
		s.ConfidenceScore,
		pq.Array(s.ValidationErrors),
		resultJSON,
	)
	return err
}

func (r *refinementRepository) CreateIteration(ctx context.Context, i *domain.RefinementIteration) error {
	query := `
		INSERT INTO refinement_iterations
		(id, session_id, iteration, prompt, response, artifact, validation_result, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	artifactJSON, _ := json.Marshal(i.Artifact)
	validationJSON, _ := json.Marshal(i.ValidationResult)

	_, err := r.db.ExecContext(ctx, query,
		i.ID,
		i.SessionID,
		i.Iteration,
		i.Prompt,
		i.Response,
		artifactJSON,
		validationJSON,
		i.CreatedAt,
	)
	return err
}
