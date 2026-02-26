package infra

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type featureIntelligenceRepository struct {
	queries *db.Queries
}

func NewFeatureIntelligenceRepository(queries *db.Queries) app.FeatureIntelligenceRepository {
	return &featureIntelligenceRepository{queries: queries}
}

func (r *featureIntelligenceRepository) Get(ctx context.Context, featureID uuid.UUID) (*domain.FeatureIntelligence, error) {
	row, err := r.queries.GetFeatureIntelligence(ctx, featureID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if not found, let service handle it
		}
		return nil, err
	}
	return &domain.FeatureIntelligence{
		ID:                       row.ID,
		FeatureID:                row.FeatureID,
		CompletenessScore:        int(row.CompletenessScore),
		ContractIntegrityScore:   int(row.ContractIntegrityScore),
		VariableCoverageScore:    int(row.VariableCoverageScore),
		DependencyStabilityScore: int(row.DependencyStabilityScore),
		DriftRiskScore:           int(row.DriftRiskScore),
		TestCoverageScore:        int(row.TestCoverageScore),
		LLMConfidenceScore:       int(row.LlmConfidenceScore),
		OverallScore:             int(row.OverallScore),
		LastCalculatedAt:         row.LastCalculatedAt.Time,
	}, nil
}

func (r *featureIntelligenceRepository) Create(ctx context.Context, fi *domain.FeatureIntelligence) error {
	_, err := r.queries.CreateFeatureIntelligence(ctx, db.CreateFeatureIntelligenceParams{
		FeatureID:                fi.FeatureID,
		CompletenessScore:        int32(fi.CompletenessScore),
		ContractIntegrityScore:   int32(fi.ContractIntegrityScore),
		VariableCoverageScore:    int32(fi.VariableCoverageScore),
		DependencyStabilityScore: int32(fi.DependencyStabilityScore),
		DriftRiskScore:           int32(fi.DriftRiskScore),
		TestCoverageScore:        int32(fi.TestCoverageScore),
		LlmConfidenceScore:       int32(fi.LLMConfidenceScore),
		OverallScore:             int32(fi.OverallScore),
	})
	return err
}

func (r *featureIntelligenceRepository) Update(ctx context.Context, fi *domain.FeatureIntelligence) error {
	_, err := r.queries.UpdateFeatureIntelligence(ctx, db.UpdateFeatureIntelligenceParams{
		FeatureID:                fi.FeatureID,
		CompletenessScore:        int32(fi.CompletenessScore),
		ContractIntegrityScore:   int32(fi.ContractIntegrityScore),
		VariableCoverageScore:    int32(fi.VariableCoverageScore),
		DependencyStabilityScore: int32(fi.DependencyStabilityScore),
		DriftRiskScore:           int32(fi.DriftRiskScore),
		TestCoverageScore:        int32(fi.TestCoverageScore),
		LlmConfidenceScore:       int32(fi.LLMConfidenceScore),
		OverallScore:             int32(fi.OverallScore),
	})
	return err
}

func (r *featureIntelligenceRepository) Delete(ctx context.Context, featureID uuid.UUID) error {
	return r.queries.DeleteFeatureIntelligence(ctx, featureID)
}
