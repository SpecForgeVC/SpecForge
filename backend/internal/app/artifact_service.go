package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type ArtifactService interface {
	GenerateArtifact(ctx context.Context, roadmapItemID uuid.UUID, format domain.ExportFormat, options ExportOptions, userID uuid.UUID) (*domain.BuildArtifactPackage, error)
}

type ExportOptions struct {
	IncludeDependencies bool `json:"include_dependencies"`
	IncludeGovernance   bool `json:"include_governance"`
}

type buildArtifactService struct {
	roadmapRepo     RoadmapItemRepository
	contractRepo    ContractRepository
	variableRepo    VariableRepository
	requirementRepo RequirementRepository
	validationRepo  ValidationRuleRepository
	govService      GovernanceService
}

func NewBuildArtifactService(
	roadmapRepo RoadmapItemRepository,
	contractRepo ContractRepository,
	variableRepo VariableRepository,
	requirementRepo RequirementRepository,
	validationRepo ValidationRuleRepository,
	govService GovernanceService,
) ArtifactService {
	return &buildArtifactService{
		roadmapRepo:     roadmapRepo,
		contractRepo:    contractRepo,
		variableRepo:    variableRepo,
		requirementRepo: requirementRepo,
		validationRepo:  validationRepo,
		govService:      govService,
	}
}

func (s *buildArtifactService) GenerateArtifact(ctx context.Context, roadmapItemID uuid.UUID, format domain.ExportFormat, options ExportOptions, userID uuid.UUID) (*domain.BuildArtifactPackage, error) {
	// 1. Fetch Roadmap Item
	item, err := s.roadmapRepo.Get(ctx, roadmapItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roadmap item: %w", err)
	}

	// 2. Build Context
	roadmapContext := domain.RoadmapContext{
		Title:            item.Title,
		Description:      item.Description,
		BusinessContext:  item.BusinessContext,
		TechnicalContext: item.TechnicalContext,
		Priority:         string(item.Priority),
		RiskLevel:        string(item.RiskLevel),
	}

	// 3. Fetch Contracts
	contracts, err := s.contractRepo.List(ctx, roadmapItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to list contracts: %w", err)
	}

	contractBundles := make([]domain.ContractBundle, 0, len(contracts))
	for _, c := range contracts {
		contractBundles = append(contractBundles, domain.ContractBundle{
			ID:           c.ID,
			Type:         string(c.ContractType),
			Version:      c.Version,
			InputSchema:  c.InputSchema,
			OutputSchema: c.OutputSchema,
		})
	}

	// 4. Fetch Requirements (Acceptance Criteria)
	reqs, err := s.requirementRepo.List(ctx, roadmapItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to list requirements: %w", err)
	}

	acceptanceCriteria := make([]domain.AcceptanceCriteria, 0, len(reqs))
	for _, r := range reqs {
		acceptanceCriteria = append(acceptanceCriteria, domain.AcceptanceCriteria{
			ID:          r.ID.String(),
			Description: r.AcceptanceCriteria,
		})
	}

	// 5. Fetch Variables for each contract
	variableBundles := make([]domain.VariableBundle, 0)
	for _, c := range contracts {
		vars, err := s.variableRepo.List(ctx, c.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list variables for contract %s: %w", c.ID, err)
		}
		for _, v := range vars {
			variableBundles = append(variableBundles, domain.VariableBundle{
				Name:            v.Name,
				Type:            v.Type,
				Required:        v.Required,
				ValidationRules: v.ValidationRules,
			})
		}
	}

	// 6. Fetch Validation Rules
	rules, err := s.validationRepo.List(ctx, item.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list validation rules: %w", err)
	}

	validationBundles := make([]domain.ValidationBundle, 0, len(rules))
	for _, r := range rules {
		validationBundles = append(validationBundles, domain.ValidationBundle{
			Name:     r.Name,
			RuleType: r.RuleType,
			Config:   r.RuleConfig,
		})
	}

	// 7. Governance Check
	var govBundle domain.GovernanceBundle
	if options.IncludeGovernance {
		ready, issues, _ := s.govService.CanBuildFeature(ctx, roadmapItemID)
		govBundle.PoliciesEnforced = issues
		if ready {
			govBundle.ComplianceStatus = "READY"
		} else {
			govBundle.ComplianceStatus = "NEEDS_REFINEMENT"
		}
	}

	// 8. Generate Prompts
	buildPrompts := s.generateBuildPrompts(roadmapContext, contractBundles, variableBundles, validationBundles, acceptanceCriteria)
	refinementPrompts := s.generateRefinementPrompts(roadmapContext, contractBundles)

	// 9. Final Package
	pkg := &domain.BuildArtifactPackage{
		Metadata: domain.MetadataSection{
			ArtifactID:     uuid.New(),
			RoadmapItemID:  roadmapItemID,
			Version:        "1.0.0",
			ExportedAt:     time.Now(),
			ExportedBy:     userID,
			IntegrityHash:  "SHA256:TODO-HASH", // TODO: Implement actual hashing
			GovernanceMode: string(item.Status),
		},
		RoadmapContext:        roadmapContext,
		Contracts:             contractBundles,
		ValidationRules:       validationBundles,
		Variables:             variableBundles,
		AcceptanceCriteria:    acceptanceCriteria,
		BuildPrompts:          buildPrompts,
		RefinementLoopPrompts: refinementPrompts,
		GovernanceConstraints: govBundle,
		Dependencies: domain.DependencyGraph{
			Nodes: []string{item.Title}, // Simple for now
			Edges: []domain.DependencyEdge{},
		},
	}

	return pkg, nil
}

func (s *buildArtifactService) generateBuildPrompts(
	ctx domain.RoadmapContext,
	contracts []domain.ContractBundle,
	vars []domain.VariableBundle,
	rules []domain.ValidationBundle,
	ac []domain.AcceptanceCriteria,
) domain.BuildPromptBundle {

	prompt := fmt.Sprintf("# IMPLEMENTATION PLAN: %s\n\n", ctx.Title)
	prompt += fmt.Sprintf("## OBJECTIVE\n%s\n\n", ctx.Description)
	prompt += "## BUSINESS CONTEXT\n" + ctx.BusinessContext + "\n\n"
	prompt += "## TECHNICAL CONTEXT\n" + ctx.TechnicalContext + "\n\n"

	prompt += "## CONTRACTS\n"
	for _, c := range contracts {
		prompt += fmt.Sprintf("- Type: %s, Version: %s\n", c.Type, c.Version)
	}
	prompt += "\n"

	prompt += "## ACCEPTANCE CRITERIA\n"
	for _, criteria := range ac {
		prompt += fmt.Sprintf("- [ ] %s\n", criteria.Description)
	}

	verification := "## VERIFICATION INSTRUCTIONS\n1. Verify implementation against schemas.\n2. Run validation rules.\n3. Ensure AC are met."

	return domain.BuildPromptBundle{
		Implementation: prompt,
		Verification:   verification,
	}
}

func (s *buildArtifactService) generateRefinementPrompts(ctx domain.RoadmapContext, contracts []domain.ContractBundle) domain.RefinementLoopBundle {
	return domain.RefinementLoopBundle{
		Instructions: "Compare output types to contract definitions. Ensure all validation rules are enforced. suggest refactor.",
	}
}
