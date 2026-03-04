package app

import (
	"context"
	"fmt"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
	"github.com/google/uuid"
)

type governanceService struct {
	fiService    FeatureIntelligenceService
	propRepo     AiProposalRepository
	varRepo      VariableRepository
	contractRepo ContractRepository
}

func NewGovernanceService(fi FeatureIntelligenceService, propRepo AiProposalRepository, varRepo VariableRepository, contractRepo ContractRepository) GovernanceService {
	return &governanceService{
		fiService:    fi,
		propRepo:     propRepo,
		varRepo:      varRepo,
		contractRepo: contractRepo,
	}
}

func (s *governanceService) CanBuildFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error) {
	var reasons []string
	allowed := true

	// 1. Check Intelligence Score
	score, err := s.fiService.GetFeatureScore(ctx, featureID)
	if err != nil {
		return false, nil, err
	}
	if score != nil && score.OverallScore < 50 {
		allowed = false
		reasons = append(reasons, fmt.Sprintf("Overall Intelligence Score is too low (%d < 50)", score.OverallScore))
	}

	// 2. Check for pending AI proposals on this feature
	proposals, err := s.propRepo.ListByRoadmapItem(ctx, featureID)
	if err != nil {
		return false, nil, err
	}
	pendingCount := 0
	for _, p := range proposals {
		if p.Status == domain.Pending {
			pendingCount++
		}
	}
	if pendingCount > 0 {
		allowed = false
		reasons = append(reasons, fmt.Sprintf("Feature has %d pending AI proposal(s) awaiting review before it can be built", pendingCount))
	}

	return allowed, reasons, nil
}

func (s *governanceService) CanDeployFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error) {
	var reasons []string
	allowed := true

	score, err := s.fiService.GetFeatureScore(ctx, featureID)
	if err != nil {
		return false, nil, err
	}
	if score != nil && score.OverallScore < 80 {
		allowed = false
		reasons = append(reasons, fmt.Sprintf("Overall Intelligence Score is too low for deployment (%d < 80)", score.OverallScore))
	}
	if score != nil && score.CompletenessScore < 100 {
		allowed = false
		reasons = append(reasons, "Feature Spec must be 100% complete before deployment")
	}

	// Also block deployment if pending proposals exist
	proposals, err := s.propRepo.ListByRoadmapItem(ctx, featureID)
	if err != nil {
		return false, nil, err
	}
	for _, p := range proposals {
		if p.Status == domain.Pending {
			allowed = false
			reasons = append(reasons, "All pending AI proposals must be reviewed before deployment")
			break
		}
	}

	return allowed, reasons, nil
}

func (s *governanceService) CanUpdateContract(ctx context.Context, contractID uuid.UUID) (bool, []string, error) {
	var reasons []string
	allowed := true

	// Resolve which roadmap item this contract belongs to
	contract, err := s.contractRepo.Get(ctx, contractID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to resolve contract: %w", err)
	}

	// Check for pending proposals on the parent roadmap item
	proposals, err := s.propRepo.ListByRoadmapItem(ctx, contract.RoadmapItemID)
	if err != nil {
		return false, nil, err
	}
	for _, p := range proposals {
		if p.Status == domain.Pending {
			allowed = false
			reasons = append(reasons,
				fmt.Sprintf("Contract has a pending AI proposal (ID: %s) awaiting review — resolve it before making manual changes", p.ID))
			break
		}
	}

	return allowed, reasons, nil
}
