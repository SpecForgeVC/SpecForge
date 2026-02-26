package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type bootstrapRepository struct {
	queries *db.Queries
}

func NewBootstrapRepository(queries *db.Queries) app.BootstrapRepository {
	return &bootstrapRepository{queries: queries}
}

// --- Snapshot mapping ---

func (r *bootstrapRepository) snapshotToDomain(row db.ProjectIntelligenceSnapshot) *domain.ProjectIntelligenceSnapshot {
	archScore, _ := strconv.ParseFloat(row.ArchitectureScore.String, 64)
	contractDensity, _ := strconv.ParseFloat(row.ContractDensity.String, 64)
	riskScore, _ := strconv.ParseFloat(row.RiskScore.String, 64)
	alignScore, _ := strconv.ParseFloat(row.AlignmentScore.String, 64)

	return &domain.ProjectIntelligenceSnapshot{
		ID:                row.ID,
		ProjectID:         row.ProjectID,
		Version:           int(row.Version),
		SnapshotJSON:      db.RawMessageToJSON(row.SnapshotJson),
		ArchitectureScore: archScore,
		ContractDensity:   contractDensity,
		RiskScore:         riskScore,
		AlignmentScore:    alignScore,
		ConfidenceJSON:    db.SqlToJSON(row.ConfidenceJson),
		CreatedAt:         row.CreatedAt,
	}
}

// --- CRUD Methods ---

func (r *bootstrapRepository) InsertSnapshot(ctx context.Context, s *domain.ProjectIntelligenceSnapshot) error {
	snapshotBytes, err := json.Marshal(s.SnapshotJSON)
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	row, err := r.queries.InsertIntelligenceSnapshot(ctx, db.InsertIntelligenceSnapshotParams{
		ProjectID:         s.ProjectID,
		Version:           int32(s.Version),
		SnapshotJson:      snapshotBytes,
		ArchitectureScore: db.TextToSql(strconv.FormatFloat(s.ArchitectureScore, 'f', 2, 64)),
		ContractDensity:   db.TextToSql(strconv.FormatFloat(s.ContractDensity, 'f', 2, 64)),
		RiskScore:         db.TextToSql(strconv.FormatFloat(s.RiskScore, 'f', 2, 64)),
		AlignmentScore:    db.TextToSql(strconv.FormatFloat(s.AlignmentScore, 'f', 2, 64)),
		ConfidenceJson:    db.JSONToSql(s.ConfidenceJSON),
	})
	if err != nil {
		return err
	}
	s.ID = row.ID
	s.CreatedAt = row.CreatedAt
	return nil
}

func (r *bootstrapRepository) ListSnapshots(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectIntelligenceSnapshot, error) {
	rows, err := r.queries.ListSnapshotsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ProjectIntelligenceSnapshot, len(rows))
	for i, row := range rows {
		result[i] = *r.snapshotToDomain(row)
	}
	return result, nil
}

func (r *bootstrapRepository) GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error) {
	row, err := r.queries.GetIntelligenceSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.snapshotToDomain(row), nil
}

func (r *bootstrapRepository) GetLatestSnapshot(ctx context.Context, projectID uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error) {
	row, err := r.queries.GetLatestSnapshot(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return r.snapshotToDomain(row), nil
}

func (r *bootstrapRepository) GetMaxVersion(ctx context.Context, projectID uuid.UUID) (int, error) {
	v, err := r.queries.GetMaxSnapshotVersion(ctx, projectID)
	return int(v), err
}

// --- Modules ---

func (r *bootstrapRepository) InsertModule(ctx context.Context, m *domain.ProjectModule) error {
	row, err := r.queries.InsertProjectModule(ctx, db.InsertProjectModuleParams{
		ProjectID:         m.ProjectID,
		SnapshotID:        m.SnapshotID,
		Name:              m.Name,
		Description:       db.TextToSql(m.Description),
		RiskLevel:         db.TextToSql(m.RiskLevel),
		ChangeSensitivity: db.TextToSql(m.ChangeSensitivity),
	})
	if err != nil {
		return err
	}
	m.ID = row.ID
	return nil
}

func (r *bootstrapRepository) ListModulesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectModule, error) {
	rows, err := r.queries.ListModulesBySnapshot(ctx, snapshotID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ProjectModule, len(rows))
	for i, row := range rows {
		result[i] = domain.ProjectModule{
			ID:                row.ID,
			ProjectID:         row.ProjectID,
			SnapshotID:        row.SnapshotID,
			Name:              row.Name,
			Description:       row.Description.String,
			RiskLevel:         row.RiskLevel.String,
			ChangeSensitivity: row.ChangeSensitivity.String,
		}
	}
	return result, nil
}

// --- Entities ---

func (r *bootstrapRepository) InsertEntity(ctx context.Context, e *domain.ProjectEntity) error {
	row, err := r.queries.InsertProjectEntity(ctx, db.InsertProjectEntityParams{
		ProjectID:         e.ProjectID,
		SnapshotID:        e.SnapshotID,
		Name:              e.Name,
		RelationshipsJson: db.JSONToSql(e.RelationshipsJSON),
		ConstraintsJson:   db.JSONToSql(e.ConstraintsJSON),
	})
	if err != nil {
		return err
	}
	e.ID = row.ID
	return nil
}

func (r *bootstrapRepository) ListEntitiesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectEntity, error) {
	rows, err := r.queries.ListEntitiesBySnapshot(ctx, snapshotID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ProjectEntity, len(rows))
	for i, row := range rows {
		result[i] = domain.ProjectEntity{
			ID:                row.ID,
			ProjectID:         row.ProjectID,
			SnapshotID:        row.SnapshotID,
			Name:              row.Name,
			RelationshipsJSON: db.SqlToJSON(row.RelationshipsJson),
			ConstraintsJSON:   db.SqlToJSON(row.ConstraintsJson),
		}
	}
	return result, nil
}

// --- API Entries ---

func (r *bootstrapRepository) InsertApiEntry(ctx context.Context, a *domain.ProjectApiEntry) error {
	row, err := r.queries.InsertProjectApiEntry(ctx, db.InsertProjectApiEntryParams{
		ProjectID:      a.ProjectID,
		SnapshotID:     a.SnapshotID,
		Endpoint:       a.Endpoint,
		Method:         a.Method,
		AuthType:       db.TextToSql(a.AuthType),
		RequestSchema:  db.JSONToSql(a.RequestSchema),
		ResponseSchema: db.JSONToSql(a.ResponseSchema),
	})
	if err != nil {
		return err
	}
	a.ID = row.ID
	return nil
}

func (r *bootstrapRepository) ListApiEntriesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectApiEntry, error) {
	rows, err := r.queries.ListApiEntriesBySnapshot(ctx, snapshotID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ProjectApiEntry, len(rows))
	for i, row := range rows {
		result[i] = domain.ProjectApiEntry{
			ID:             row.ID,
			ProjectID:      row.ProjectID,
			SnapshotID:     row.SnapshotID,
			Endpoint:       row.Endpoint,
			Method:         row.Method,
			AuthType:       row.AuthType.String,
			RequestSchema:  db.SqlToJSON(row.RequestSchema),
			ResponseSchema: db.SqlToJSON(row.ResponseSchema),
		}
	}
	return result, nil
}

// --- Contract Entries ---

func (r *bootstrapRepository) InsertContractEntry(ctx context.Context, c *domain.ProjectContractEntry) error {
	row, err := r.queries.InsertProjectContractEntry(ctx, db.InsertProjectContractEntryParams{
		ProjectID:      c.ProjectID,
		SnapshotID:     c.SnapshotID,
		Name:           c.Name,
		ContractType:   db.TextToSql(c.ContractType),
		SchemaJson:     db.JSONToSql(c.SchemaJSON),
		SourceModule:   db.TextToSql(c.SourceModule),
		StabilityScore: db.TextToSql(strconv.FormatFloat(c.StabilityScore, 'f', 2, 64)),
	})
	if err != nil {
		return err
	}
	c.ID = row.ID
	return nil
}

func (r *bootstrapRepository) ListContractEntriesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectContractEntry, error) {
	rows, err := r.queries.ListContractEntriesBySnapshot(ctx, snapshotID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ProjectContractEntry, len(rows))
	for i, row := range rows {
		stabilityScore, _ := strconv.ParseFloat(row.StabilityScore.String, 64)
		result[i] = domain.ProjectContractEntry{
			ID:             row.ID,
			ProjectID:      row.ProjectID,
			SnapshotID:     row.SnapshotID,
			Name:           row.Name,
			ContractType:   row.ContractType.String,
			SchemaJSON:     db.SqlToJSON(row.SchemaJson),
			SourceModule:   row.SourceModule.String,
			StabilityScore: stabilityScore,
		}
	}
	return result, nil
}
