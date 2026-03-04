package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/SpecForgeVC/SpecForge/internal/domain"
	"github.com/google/uuid"
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
	fiService       FeatureIntelligenceService
}

func NewBuildArtifactService(
	roadmapRepo RoadmapItemRepository,
	contractRepo ContractRepository,
	variableRepo VariableRepository,
	requirementRepo RequirementRepository,
	validationRepo ValidationRuleRepository,
	govService GovernanceService,
	fiService FeatureIntelligenceService,
) ArtifactService {
	return &buildArtifactService{
		roadmapRepo:     roadmapRepo,
		contractRepo:    contractRepo,
		variableRepo:    variableRepo,
		requirementRepo: requirementRepo,
		validationRepo:  validationRepo,
		govService:      govService,
		fiService:       fiService,
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

	// 2a. Fetch readiness score from Feature Intelligence
	if fi, err := s.fiService.GetFeatureScore(ctx, roadmapItemID); err == nil && fi != nil {
		roadmapContext.ReadinessScore = fi.OverallScore
		roadmapContext.ReadinessLevel = string(item.ReadinessLevel)
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

	// 8. Build the dependency graph from actual contract/requirement names
	var depNodes []string
	var depEdges []domain.DependencyEdge
	depNodes = append(depNodes, item.Title)
	for _, c := range contracts {
		contractLabel := fmt.Sprintf("%s v%s", c.ContractType, c.Version)
		depNodes = append(depNodes, contractLabel)
		depEdges = append(depEdges, domain.DependencyEdge{From: item.Title, To: contractLabel, Type: "contract"})
	}
	for _, r := range reqs {
		depNodes = append(depNodes, r.Title)
		depEdges = append(depEdges, domain.DependencyEdge{From: item.Title, To: r.Title, Type: "requirement"})
	}

	// 9. Generate Prompts
	buildPrompts := s.generateBuildPrompts(roadmapContext, contractBundles, variableBundles, validationBundles, acceptanceCriteria)
	refinementPrompts := s.generateRefinementPrompts(roadmapContext, contractBundles)

	// 10. Compute integrity hash over the core package data
	integrityHash := computeIntegrityHash(roadmapContext, contractBundles, variableBundles, acceptanceCriteria)

	// 11. Final Package
	pkg := &domain.BuildArtifactPackage{
		Metadata: domain.MetadataSection{
			ArtifactID:     uuid.New(),
			RoadmapItemID:  roadmapItemID,
			Version:        "1.0.0",
			ExportedAt:     time.Now(),
			ExportedBy:     userID,
			IntegrityHash:  integrityHash,
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
			Nodes: depNodes,
			Edges: depEdges,
		},
	}

	return pkg, nil
}

// computeIntegrityHash produces a SHA-256 hex digest over key package fields.
func computeIntegrityHash(ctx domain.RoadmapContext, contracts []domain.ContractBundle, vars []domain.VariableBundle, ac []domain.AcceptanceCriteria) string {
	h := sha256.New()
	data, _ := json.Marshal(map[string]interface{}{
		"context":   ctx,
		"contracts": contracts,
		"variables": vars,
		"ac":        ac,
	})
	h.Write(data)
	return "SHA256:" + hex.EncodeToString(h.Sum(nil))
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
	prompt += fmt.Sprintf("## RISK & PRIORITY\n- Priority: **%s**\n- Risk Level: **%s**\n- Readiness Score: **%d%%** (%s)\n\n",
		ctx.Priority, ctx.RiskLevel, ctx.ReadinessScore, ctx.ReadinessLevel)

	prompt += "## CONTRACTS\n"
	if len(contracts) == 0 {
		prompt += "_No contracts defined for this feature._\n"
	} else {
		for _, c := range contracts {
			prompt += fmt.Sprintf("- **Type**: %s | **Version**: %s | **ID**: %s\n", c.Type, c.Version, c.ID)
		}
	}
	prompt += "\n"

	if len(vars) > 0 {
		prompt += "## VARIABLES\n"
		for _, v := range vars {
			required := ""
			if v.Required {
				required = " _(required)_"
			}
			prompt += fmt.Sprintf("- `%s` (%s)%s\n", v.Name, v.Type, required)
		}
		prompt += "\n"
	}

	if len(rules) > 0 {
		prompt += "## VALIDATION RULES\n"
		for _, r := range rules {
			prompt += fmt.Sprintf("- **%s** (%s)\n", r.Name, r.RuleType)
		}
		prompt += "\n"
	}

	prompt += "## ACCEPTANCE CRITERIA\n"
	if len(ac) == 0 {
		prompt += "_No acceptance criteria defined._\n"
	} else {
		for _, criteria := range ac {
			prompt += fmt.Sprintf("- [ ] %s\n", criteria.Description)
		}
	}

	contractTypes := make([]string, 0, len(contracts))
	for _, c := range contracts {
		contractTypes = append(contractTypes, c.Type)
	}
	contractSummary := "N/A"
	if len(contractTypes) > 0 {
		contractSummary = strings.Join(contractTypes, ", ")
	}

	verification := fmt.Sprintf(`## VERIFICATION INSTRUCTIONS
1. Verify all contract schemas (%s) are fully implemented.
2. Ensure all validation rules are enforced at the API boundary.
3. Run acceptance criteria checks against the implementation.
4. Readiness score must be >= 80%% before deploying (current: %d%%).
5. If risk level is HIGH or CRITICAL, request a peer review before merging.`,
		contractSummary, ctx.ReadinessScore)

	return domain.BuildPromptBundle{
		Implementation: prompt,
		Verification:   verification,
	}
}

func (s *buildArtifactService) generateRefinementPrompts(ctx domain.RoadmapContext, contracts []domain.ContractBundle) domain.RefinementLoopBundle {
	contractDescs := make([]string, 0, len(contracts))
	for _, c := range contracts {
		contractDescs = append(contractDescs, fmt.Sprintf("%s v%s", c.Type, c.Version))
	}
	contractList := "none"
	if len(contractDescs) > 0 {
		contractList = strings.Join(contractDescs, ", ")
	}
	return domain.RefinementLoopBundle{
		Instructions: fmt.Sprintf(
			"Compare output types against the defined contracts (%s). "+
				"Ensure all validation rules are enforced. "+
				"If the readiness score is below 80%%, identify and address missing contracts, variables, or acceptance criteria. "+
				"Suggest a targeted refactor to close any gaps.",
			contractList,
		),
	}
}
