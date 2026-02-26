package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra/db"
)

type aiProposalRepository struct {
	queries *db.Queries
}

func NewAiProposalRepository(queries *db.Queries) app.AiProposalRepository {
	return &aiProposalRepository{queries: queries}
}

func (r *aiProposalRepository) Get(ctx context.Context, id uuid.UUID) (*domain.AiProposal, error) {
	row, err := r.queries.GetAiProposal(ctx, id)
	if err != nil {
		return nil, err
	}
	var reviewedBy uuid.UUID
	if row.ReviewedBy.Valid {
		reviewedBy = row.ReviewedBy.UUID
	}

	return &domain.AiProposal{
		ID:              row.ID,
		RoadmapItemID:   row.RoadmapItemID,
		ProposalType:    domain.ProposalType(row.ProposalType),
		Diff:            db.RawMessageToJSON(row.Diff),
		Reasoning:       row.Reasoning.String,
		ConfidenceScore: row.ConfidenceScore.Float64,
		Status:          domain.ProposalStatus(row.Status.ProposalStatus),
		ReviewedBy:      reviewedBy,
		CreatedAt:       row.CreatedAt.Time,
	}, nil
}

func (r *aiProposalRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.AiProposal, error) {
	rows, err := r.queries.ListAiProposalsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return r.mapRows(rows), nil
}

func (r *aiProposalRepository) ListByRoadmapItem(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.AiProposal, error) {
	rows, err := r.queries.ListAiProposalsByRoadmapItem(ctx, roadmapItemID)
	if err != nil {
		return nil, err
	}
	return r.mapRows(rows), nil
}

func (r *aiProposalRepository) mapRows(rows []db.AiProposal) []domain.AiProposal {
	proposals := make([]domain.AiProposal, len(rows))
	for i, row := range rows {
		var reviewedBy uuid.UUID
		if row.ReviewedBy.Valid {
			reviewedBy = row.ReviewedBy.UUID
		}
		proposals[i] = domain.AiProposal{
			ID:              row.ID,
			RoadmapItemID:   row.RoadmapItemID,
			ProposalType:    domain.ProposalType(row.ProposalType),
			Diff:            db.RawMessageToJSON(row.Diff),
			Reasoning:       row.Reasoning.String,
			ConfidenceScore: row.ConfidenceScore.Float64,
			Status:          domain.ProposalStatus(row.Status.ProposalStatus),
			ReviewedBy:      reviewedBy,
			CreatedAt:       row.CreatedAt.Time,
		}
	}
	return proposals
}

func (r *aiProposalRepository) Create(ctx context.Context, p *domain.AiProposal) error {
	_, err := r.queries.CreateAiProposal(ctx, db.CreateAiProposalParams{
		RoadmapItemID:   p.RoadmapItemID,
		ProposalType:    db.ProposalType(p.ProposalType),
		Diff:            db.JSONToRawMessage(p.Diff),
		Reasoning:       db.TextToSql(p.Reasoning),
		ConfidenceScore: db.FloatToSql(p.ConfidenceScore),
	})
	return err
}

func (r *aiProposalRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProposalStatus, reviewedBy uuid.UUID) error {
	_, err := r.queries.UpdateAiProposalStatus(ctx, db.UpdateAiProposalStatusParams{
		ID:         id,
		Status:     db.NullProposalStatus{ProposalStatus: db.ProposalStatus(status), Valid: true},
		ReviewedBy: uuid.NullUUID{UUID: reviewedBy, Valid: true},
	})
	return err
}
