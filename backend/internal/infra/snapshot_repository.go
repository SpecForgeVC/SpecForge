package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type snapshotRepository struct {
	queries *db.Queries
}

func NewSnapshotRepository(queries *db.Queries) app.SnapshotRepository {
	return &snapshotRepository{queries: queries}
}

func (r *snapshotRepository) Get(ctx context.Context, id uuid.UUID) (*domain.VersionSnapshot, error) {
	row, err := r.queries.GetVersionSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}
	var createdBy uuid.UUID
	if row.CreatedBy.Valid {
		createdBy = row.CreatedBy.UUID
	}

	return &domain.VersionSnapshot{
		ID:            row.ID,
		RoadmapItemID: row.RoadmapItemID,
		SnapshotData:  db.RawMessageToJSON(row.SnapshotData),
		CreatedAt:     row.CreatedAt.Time,
		CreatedBy:     createdBy,
	}, nil
}

func (r *snapshotRepository) List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.VersionSnapshot, error) {
	rows, err := r.queries.ListVersionSnapshots(ctx, roadmapItemID)
	if err != nil {
		return nil, err
	}
	snapshots := make([]domain.VersionSnapshot, len(rows))
	for i, row := range rows {
		snapshots[i] = *r.mapRow(row)
	}
	return snapshots, nil
}

func (r *snapshotRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VersionSnapshot, error) {
	rows, err := r.queries.ListVersionSnapshotsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	snapshots := make([]domain.VersionSnapshot, len(rows))
	for i, row := range rows {
		snapshots[i] = *r.mapRow(row)
	}
	return snapshots, nil
}

func (r *snapshotRepository) mapRow(row db.VersionSnapshot) *domain.VersionSnapshot {
	var createdBy uuid.UUID
	if row.CreatedBy.Valid {
		createdBy = row.CreatedBy.UUID
	}

	return &domain.VersionSnapshot{
		ID:            row.ID,
		RoadmapItemID: row.RoadmapItemID,
		SnapshotData:  db.RawMessageToJSON(row.SnapshotData),
		CreatedAt:     row.CreatedAt.Time,
		CreatedBy:     createdBy,
	}
}

func (r *snapshotRepository) Create(ctx context.Context, s *domain.VersionSnapshot) error {
	row, err := r.queries.CreateVersionSnapshot(ctx, db.CreateVersionSnapshotParams{
		RoadmapItemID: s.RoadmapItemID,
		SnapshotData:  db.JSONToRawMessage(s.SnapshotData),
		CreatedBy:     uuid.NullUUID{UUID: s.CreatedBy, Valid: true},
	})
	if err != nil {
		return err
	}
	s.ID = row.ID
	if row.CreatedAt.Valid {
		s.CreatedAt = row.CreatedAt.Time
	}
	return nil
}
