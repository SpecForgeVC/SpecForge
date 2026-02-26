package infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type MCPRepository struct {
	db *sql.DB
}

func NewMCPRepository(db *sql.DB) *MCPRepository {
	return &MCPRepository{db: db}
}

func (r *MCPRepository) CreateSnapshot(ctx context.Context, projectID, roadmapItemID uuid.UUID) (uuid.UUID, error) {
	id := uuid.New()
	query := `
		INSERT INTO reality_snapshots (id, project_id, roadmap_item_id, state)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, id, projectID, roadmapItemID, domain.StateInitiated)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create snapshot: %w", err)
	}
	return id, nil
}

func (r *MCPRepository) UpdateSnapshotState(ctx context.Context, id uuid.UUID, state domain.SnapshotState) error {
	query := `
		UPDATE reality_snapshots
		SET state = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, state, id)
	return err
}

func (r *MCPRepository) SaveSnapshotData(ctx context.Context, id uuid.UUID, data domain.EnvironmentSnapshot) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	query := `
		UPDATE reality_snapshots
		SET snapshot_json = $1, state = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	_, err = r.db.ExecContext(ctx, query, jsonData, domain.StateAnalyzing, id)
	return err
}

func (r *MCPRepository) GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.SnapshotData, error) {
	query := `
		SELECT id, project_id, roadmap_item_id, state, snapshot_json, created_at, updated_at
		FROM reality_snapshots
		WHERE id = $1
	`
	var snapshot domain.SnapshotData
	var jsonData []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&snapshot.ID, &snapshot.ProjectID, &snapshot.RoadmapItemID, &snapshot.State, &jsonData, &snapshot.CreatedAt, &snapshot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if jsonData != nil {
		if err := json.Unmarshal(jsonData, &snapshot.SnapshotJSON); err != nil {
			return nil, err
		}
	}

	return &snapshot, nil
}

func (r *MCPRepository) SaveAnalysis(ctx context.Context, analysis domain.SnapshotAnalysis) error {
	scoresJSON, err := json.Marshal(analysis.Scores)
	if err != nil {
		return err
	}
	conflictsJSON, err := json.Marshal(analysis.AlignmentConflicts)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO snapshot_analysis (snapshot_id, scores_json, verdict, drift_detected, alignment_conflicts)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.ExecContext(ctx, query, analysis.SnapshotID, scoresJSON, analysis.Verdict, analysis.DriftDetected, conflictsJSON)
	return err
}

func (r *MCPRepository) ListActiveSnapshots(ctx context.Context, projectID uuid.UUID) ([]domain.SnapshotData, error) {
	query := `
		SELECT id, project_id, roadmap_item_id, state, created_at, updated_at
		FROM reality_snapshots
		WHERE project_id = $1 AND state NOT IN ('completed', 'failed')
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []domain.SnapshotData
	for rows.Next() {
		var s domain.SnapshotData
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.RoadmapItemID, &s.State, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}
