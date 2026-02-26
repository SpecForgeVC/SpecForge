package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type governanceService struct {
	fiService FeatureIntelligenceService
	propRepo  AiProposalRepository
	varRepo   VariableRepository
}

func NewGovernanceService(fi FeatureIntelligenceService, propRepo AiProposalRepository, varRepo VariableRepository) GovernanceService {
	return &governanceService{
		fiService: fi,
		propRepo:  propRepo,
		varRepo:   varRepo,
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
	if score != nil && score.OverallScore < 50 { // Threshold
		allowed = false
		reasons = append(reasons, fmt.Sprintf("Overall Intelligence Score is too low (%d < 50)", score.OverallScore))
	}

	// 2. Check for Pending Proposals (Simplified check - would need more complex query in reality)
	// For now, we assume if there are ANY pending proposals for this feature's contracts, we block.
	// This requires querying proposals by roadmap item or traversing contracts.
	// Let's assume we skip this deep check for now or add a method to propRepo.

	return allowed, reasons, nil
}

func (s *governanceService) CanDeployFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error) {
	// Stricter rules for deployment
	var reasons []string
	allowed := true

	score, err := s.fiService.GetFeatureScore(ctx, featureID)
	if err != nil {
		return false, nil, err
	}
	if score != nil && score.OverallScore < 80 { // Higher Threshold
		allowed = false
		reasons = append(reasons, fmt.Sprintf("Overall Intelligence Score is too low for deployment (%d < 80)", score.OverallScore))
	}

	if score != nil && score.CompletenessScore < 100 {
		allowed = false
		reasons = append(reasons, "Feature Spec must be 100% complete")
	}

	return allowed, reasons, nil
}

func (s *governanceService) CanUpdateContract(ctx context.Context, contractID uuid.UUID) (bool, []string, error) {
	// For now, allow all updates. Real implementation would check proposals.
	return true, nil, nil
}
