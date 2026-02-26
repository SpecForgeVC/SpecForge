package infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type alignmentRepository struct {
	db *sql.DB
}

func NewAlignmentRepository(db *sql.DB) app.AlignmentRepository {
	return &alignmentRepository{db: db}
}

func (r *alignmentRepository) GetLatestReport(ctx context.Context, projectID uuid.UUID) (*domain.AlignmentReport, error) {
	query := `
		SELECT id, project_id, conflicts, "overlaps", missing_dependencies, circular_dependencies, recommended_resolutions, alignment_score, created_at
		FROM alignment_reports
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, projectID)

	var report domain.AlignmentReport
	var conflictsJSON, overlapsJSON, missingDepsJSON, circularDepsJSON, recommendedResolutionsJSON []byte

	err := row.Scan(
		&report.ID,
		&report.ProjectID,
		&conflictsJSON,
		&overlapsJSON,
		&missingDepsJSON,
		&circularDepsJSON,
		&recommendedResolutionsJSON,
		&report.AlignmentScore,
		&report.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or return a default report?
		}
		return nil, err
	}

	json.Unmarshal(conflictsJSON, &report.Conflicts)
	json.Unmarshal(overlapsJSON, &report.Overlaps)
	json.Unmarshal(missingDepsJSON, &report.MissingDependencies)
	json.Unmarshal(circularDepsJSON, &report.CircularDependencies)
	json.Unmarshal(recommendedResolutionsJSON, &report.RecommendedResolutions)

	return &report, nil
}

func (r *alignmentRepository) CreateReport(ctx context.Context, report *domain.AlignmentReport) error {
	query := `
		INSERT INTO alignment_reports (
			id, project_id, conflicts, "overlaps", missing_dependencies, circular_dependencies, recommended_resolutions, alignment_score, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	conflictsJSON, _ := json.Marshal(report.Conflicts)
	overlapsJSON, _ := json.Marshal(report.Overlaps)
	missingDepsJSON, _ := json.Marshal(report.MissingDependencies)
	circularDepsJSON, _ := json.Marshal(report.CircularDependencies)
	recommendedResolutionsJSON, _ := json.Marshal(report.RecommendedResolutions)

	if report.CreatedAt.IsZero() {
		report.CreatedAt = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query,
		report.ID,
		report.ProjectID,
		conflictsJSON,
		overlapsJSON,
		missingDepsJSON,
		circularDepsJSON,
		recommendedResolutionsJSON,
		report.AlignmentScore,
		report.CreatedAt,
	)

	// Also update project's alignment score
	if err == nil {
		updateProj := `UPDATE projects SET alignment_score = $1 WHERE id = $2`
		_, _ = r.db.ExecContext(ctx, updateProj, report.AlignmentScore, report.ProjectID)
	}

	return err
}
