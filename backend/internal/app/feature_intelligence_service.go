package app

import (
	"context"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
	"github.com/SpecForgeVC/SpecForge/internal/drift"
	"github.com/google/uuid"
)

type featureIntelligenceService struct {
	repo            FeatureIntelligenceRepository
	roadmapRepo     RoadmapItemRepository
	contractRepo    ContractRepository
	variableRepo    VariableRepository
	requirementRepo RequirementRepository
	driftService    drift.DriftService
	notifier        NotificationService
}

// Define DriftService interface locally or import if we use the package.
// Since we can't easily import internal/drift from app without cycles if app is imported by drift (not the case here, drift imports domain),
// we should import "github.com/SpecForgeVC/SpecForge/internal/drift"

// NewFeatureIntelligenceService creates a new instance of FeatureIntelligenceService
func NewFeatureIntelligenceService(
	repo FeatureIntelligenceRepository,
	roadmapRepo RoadmapItemRepository,
	contractRepo ContractRepository,
	variableRepo VariableRepository,
	requirementRepo RequirementRepository,
	driftService drift.DriftService,
	notifier NotificationService,
) FeatureIntelligenceService {
	return &featureIntelligenceService{
		repo:            repo,
		roadmapRepo:     roadmapRepo,
		contractRepo:    contractRepo,
		variableRepo:    variableRepo,
		requirementRepo: requirementRepo,
		driftService:    driftService,
		notifier:        notifier,
	}
}

func (s *featureIntelligenceService) GetFeatureScore(ctx context.Context, featureID uuid.UUID) (*domain.FeatureIntelligence, error) {
	return s.repo.Get(ctx, featureID)
}

func (s *featureIntelligenceService) CalculateFeatureScore(ctx context.Context, featureID uuid.UUID) (*domain.FeatureIntelligence, error) {
	// 1. Fetch data
	feature, err := s.roadmapRepo.Get(ctx, featureID)
	if err != nil {
		return nil, err
	}
	contracts, err := s.contractRepo.List(ctx, featureID)
	if err != nil {
		return nil, err
	}
	// Variables are per contract usually, but we might list all by walking contracts
	// For simplicity, let's assume we fetch variables for all contracts
	var totalVariables int
	for _, c := range contracts {
		vars, err := s.variableRepo.List(ctx, c.ID)
		if err == nil {
			totalVariables += len(vars)
		}
	}
	requirements, err := s.requirementRepo.List(ctx, featureID)
	if err != nil {
		return nil, err
	}

	// 2. Scoring Logic (Simplified for MVP)
	// Completeness: Fields filled in RoadmapItem
	completeness := 0
	if feature.Description != "" {
		completeness += 20
	}
	if feature.BusinessContext != "" {
		completeness += 40
	}
	if feature.TechnicalContext != "" {
		completeness += 40
	}

	// Contract Integrity: Contracts existing
	contractScore := 0
	if len(contracts) > 0 {
		contractScore = 100 // Refine later based on schema validity
	}

	// Variable Coverage: Do variables exist?
	variableScore := 0
	if totalVariables > 0 {
		variableScore = 100
	} else if len(contracts) == 0 {
		variableScore = 100 // No contracts, no variables needed logic?
	}

	// Test Coverage: Requirements testable?
	testScore := 0
	testableReqs := 0
	if len(requirements) > 0 {
		for _, r := range requirements {
			if r.Testable {
				testableReqs++
			}
		}
		testScore = (testableReqs * 100) / len(requirements)
	} else {
		testScore = 0 // No reqs = bad
	}

	// Drift Risk
	driftScore, err := s.driftService.GetFeatureDriftScore(ctx, featureID)
	if err != nil {
		// Log error but don't fail, assume 100? or 0?
		driftScore = 100
	}

	// Dependency Stability: Mock based on contract count for now, but not 100 fixed
	depScore := 100
	if len(contracts) == 0 && feature.TechnicalContext != "" {
		depScore = 50 // Expected contracts but found none
	}

	// LLM Confidence: Based on completeness of AI-relevant fields
	llmScore := (completeness + testScore) / 2
	if llmScore < 50 {
		llmScore = 50 // Base floor for AI visibility
	}

	overall := (completeness * 20 / 100) +
		(contractScore * 20 / 100) +
		(variableScore * 15 / 100) +
		(testScore * 10 / 100) +
		(depScore * 10 / 100) +
		(driftScore * 15 / 100) +
		(llmScore * 10 / 100)

	fi := &domain.FeatureIntelligence{
		FeatureID:                featureID,
		CompletenessScore:        completeness,
		ContractIntegrityScore:   contractScore,
		VariableCoverageScore:    variableScore,
		DependencyStabilityScore: depScore,
		DriftRiskScore:           driftScore,
		TestCoverageScore:        testScore,
		LLMConfidenceScore:       llmScore,
		OverallScore:             overall,
	}

	// 3. Save
	existing, err := s.repo.Get(ctx, featureID)
	if err != nil {
		// Create
		if err := s.repo.Create(ctx, fi); err != nil {
			return nil, err
		}
	} else if existing != nil {
		// Update
		if err := s.repo.Update(ctx, fi); err != nil {
			return nil, err
		}
	} else {
		// Should have been err or nil, handle create
		if err := s.repo.Create(ctx, fi); err != nil {
			return nil, err
		}
	}

	// Broadcast Update
	if s.notifier != nil {
		s.notifier.Broadcast("FEATURE_SCORE_UPDATED", fi)
	}

	return fi, nil
}
