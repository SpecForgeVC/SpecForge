package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type roadmapDependencyRepository struct {
	db *sql.DB
}

func NewRoadmapDependencyRepository(db *sql.DB) app.RoadmapDependencyRepository {
	return &roadmapDependencyRepository{db: db}
}

func (r *roadmapDependencyRepository) Create(ctx context.Context, dep *domain.RoadmapDependency) error {
	query := `
		INSERT INTO roadmap_dependencies (id, source_id, target_id, dependency_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, dep.ID, dep.SourceID, dep.TargetID, dep.DependencyType, dep.CreatedAt)
	return err
}

func (r *roadmapDependencyRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapDependency, error) {
	query := `
		SELECT d.id, d.source_id, d.target_id, d.dependency_type, d.created_at
		FROM roadmap_dependencies d
		JOIN roadmap_items r ON d.source_id = r.id
		WHERE r.project_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []domain.RoadmapDependency
	for rows.Next() {
		var d domain.RoadmapDependency
		if err := rows.Scan(&d.ID, &d.SourceID, &d.TargetID, &d.DependencyType, &d.CreatedAt); err != nil {
			return nil, err
		}
		deps = append(deps, d)
	}
	return deps, nil
}

func (r *roadmapDependencyRepository) ListBySource(ctx context.Context, sourceID uuid.UUID) ([]domain.RoadmapDependency, error) {
	query := `
		SELECT id, source_id, target_id, dependency_type, created_at
		FROM roadmap_dependencies
		WHERE source_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []domain.RoadmapDependency
	for rows.Next() {
		var d domain.RoadmapDependency
		if err := rows.Scan(&d.ID, &d.SourceID, &d.TargetID, &d.DependencyType, &d.CreatedAt); err != nil {
			return nil, err
		}
		deps = append(deps, d)
	}
	return deps, nil
}

func (r *roadmapDependencyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM roadmap_dependencies WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
