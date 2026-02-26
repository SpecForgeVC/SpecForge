package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/drift"
)

type aiProposalService struct {
	repo         AiProposalRepository
	roadmapRepo  RoadmapItemRepository
	snapshotRepo SnapshotRepository
	diffEngine   drift.DiffEngine
	auditLog     AuditLogService
}

func NewAiProposalService(repo AiProposalRepository, rmRepo RoadmapItemRepository, sRepo SnapshotRepository, de drift.DiffEngine, al AuditLogService) AiProposalService {
	return &aiProposalService{
		repo:         repo,
		roadmapRepo:  rmRepo,
		snapshotRepo: sRepo,
		diffEngine:   de,
		auditLog:     al,
	}
}

func (s *aiProposalService) GetProposal(ctx context.Context, id uuid.UUID) (*domain.AiProposal, error) {
	return s.repo.Get(ctx, id)
}

func (s *aiProposalService) ListProposals(ctx context.Context, projectID uuid.UUID) ([]domain.AiProposal, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *aiProposalService) CreateProposal(ctx context.Context, roadmapItemID uuid.UUID, pType domain.ProposalType, diff map[string]interface{}, reasoning string, confidence float64) (*domain.AiProposal, error) {
	p := &domain.AiProposal{
		ID:              uuid.New(),
		RoadmapItemID:   roadmapItemID,
		ProposalType:    pType,
		Diff:            diff,
		Reasoning:       reasoning,
		ConfidenceScore: confidence,
		Status:          domain.Pending,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	// Note: p.ID is used as entityID, roadmapItemID as context
	s.auditLog.Log(ctx, "ai_proposal", p.ID, "CREATE", uuid.Nil, nil, map[string]interface{}{"type": p.ProposalType})
	return p, nil
}

func (s *aiProposalService) ApproveProposal(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if p.Status != domain.Pending {
		return fmt.Errorf("proposal is already %s", p.Status)
	}

	// 1. Update status
	if err := s.repo.UpdateStatus(ctx, id, domain.Approved, userID); err != nil {
		return err
	}

	s.auditLog.Log(ctx, "ai_proposal", id, "APPROVE", userID, map[string]interface{}{"status": p.Status}, map[string]interface{}{"status": domain.Approved})

	// 2. Apply changes to RoadmapItem (simplified for now)
	rm, err := s.roadmapRepo.Get(ctx, p.RoadmapItemID)
	if err != nil {
		return err
	}

	// TODO: Based on ProposalType, update RoadmapItem fields
	// For now, let's just create a snapshot

	// 3. Create Version Snapshot
	snap := &domain.VersionSnapshot{
		ID:            uuid.New(),
		RoadmapItemID: p.RoadmapItemID,
		SnapshotData: map[string]interface{}{
			"proposal_id":   id,
			"proposal_type": p.ProposalType,
			"diff":          p.Diff,
			"roadmap_item":  rm,
		},
		CreatedBy: userID,
	}
	return s.snapshotRepo.Create(ctx, snap)
}

func (s *aiProposalService) RejectProposal(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateStatus(ctx, id, domain.Rejected, userID); err != nil {
		return err
	}
	s.auditLog.Log(ctx, "ai_proposal", id, "REJECT", userID, map[string]interface{}{"status": p.Status}, map[string]interface{}{"status": domain.Rejected})
	return nil
}
