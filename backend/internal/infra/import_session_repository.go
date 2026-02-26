package infra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type importSessionRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewImportSessionRepository(d *sql.DB) app.ImportSessionRepository {
	return &importSessionRepository{
		queries: db.New(d),
		db:      d,
	}
}

func (r *importSessionRepository) CreateSession(ctx context.Context, session *domain.ImportSession) error {
	res, err := r.queries.CreateImportSession(ctx, db.CreateImportSessionParams{
		ProjectID:         session.ProjectID,
		CompletenessScore: session.CompletenessScore,
		Status:            db.ImportStatus(session.Status),
		IterationCount:    session.IterationCount,
		Locked:            session.Locked,
	})
	if err != nil {
		return err
	}
	session.ID = res.ID
	session.CreatedAt = res.CreatedAt.Time
	session.UpdatedAt = res.UpdatedAt.Time
	return nil
}

func (r *importSessionRepository) GetSession(ctx context.Context, id uuid.UUID) (*domain.ImportSession, error) {
	res, err := r.queries.GetImportSession(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapSession(res), nil
}

func (r *importSessionRepository) GetLatestSessionByProject(ctx context.Context, projectID uuid.UUID) (*domain.ImportSession, error) {
	res, err := r.queries.GetLatestImportSessionByProject(ctx, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapSession(res), nil
}

func (r *importSessionRepository) UpdateSession(ctx context.Context, session *domain.ImportSession) error {
	res, err := r.queries.UpdateImportSession(ctx, db.UpdateImportSessionParams{
		ID:                session.ID,
		CompletenessScore: session.CompletenessScore,
		Status:            db.ImportStatus(session.Status),
		IterationCount:    session.IterationCount,
		Locked:            session.Locked,
	})
	if err != nil {
		return err
	}
	session.UpdatedAt = res.UpdatedAt.Time
	return nil
}

func (r *importSessionRepository) CreateArtifact(ctx context.Context, artifact *domain.ImportArtifact) error {
	payloadBytes, err := json.Marshal(artifact.Payload)
	if err != nil {
		return err
	}
	res, err := r.queries.CreateImportArtifact(ctx, db.CreateImportArtifactParams{
		SessionID: artifact.SessionID,
		Payload:   payloadBytes,
	})
	if err != nil {
		return err
	}
	artifact.ID = res.ID
	artifact.CreatedAt = res.CreatedAt.Time
	return nil
}

func (r *importSessionRepository) ListArtifactsBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.ImportArtifact, error) {
	rows, err := r.queries.ListImportArtifactsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	artifacts := make([]domain.ImportArtifact, len(rows))
	for i, row := range rows {
		var payload map[string]interface{}
		if err := json.Unmarshal(row.Payload, &payload); err != nil {
			return nil, err
		}
		artifacts[i] = domain.ImportArtifact{
			ID:        row.ID,
			SessionID: row.SessionID,
			Payload:   payload,
			CreatedAt: row.CreatedAt.Time,
		}
	}
	return artifacts, nil
}

func mapSession(row db.ProjectImportSession) *domain.ImportSession {
	return &domain.ImportSession{
		ID:                row.ID,
		ProjectID:         row.ProjectID,
		CompletenessScore: row.CompletenessScore,
		Status:            domain.ImportStatus(row.Status),
		IterationCount:    row.IterationCount,
		Locked:            row.Locked,
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
	}
}
