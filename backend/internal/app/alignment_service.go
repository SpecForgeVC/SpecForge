package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type alignmentService struct {
	repo           AlignmentRepository
	roadmapRepo    RoadmapItemRepository
	depRepo        RoadmapDependencyRepository
	contractRepo   ContractRepository
	variableRepo   VariableRepository
	validationRepo ValidationRuleRepository
}

func NewAlignmentService(
	repo AlignmentRepository,
	roadmapRepo RoadmapItemRepository,
	depRepo RoadmapDependencyRepository,
	contractRepo ContractRepository,
	variableRepo VariableRepository,
	validationRepo ValidationRuleRepository,
) AlignmentService {
	return &alignmentService{
		repo:           repo,
		roadmapRepo:    roadmapRepo,
		depRepo:        depRepo,
		contractRepo:   contractRepo,
		variableRepo:   variableRepo,
		validationRepo: validationRepo,
	}
}

func (s *alignmentService) GetAlignmentReport(ctx context.Context, projectID uuid.UUID) (*domain.AlignmentReport, error) {
	return s.repo.GetLatestReport(ctx, projectID)
}

func (s *alignmentService) TriggerAlignmentCheck(ctx context.Context, projectID uuid.UUID) (*domain.AlignmentReport, error) {
	// 1. Fetch live data from all repositories
	roadmapItems, err := s.roadmapRepo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}

	dependencies, err := s.depRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	contracts, err := s.contractRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	variables, err := s.variableRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	rules, err := s.validationRepo.List(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 2. Initialize report
	report := &domain.AlignmentReport{
		ID:                     uuid.New(),
		ProjectID:              projectID,
		Conflicts:              []domain.Conflict{},
		Overlaps:               []domain.Overlap{},
		MissingDependencies:    []string{},
		CircularDependencies:   []string{},
		RecommendedResolutions: []string{},
		AlignmentScore:         100,
		CreatedAt:              time.Now(),
	}

	// 3. Run detectors
	s.detectDependencyLoops(roadmapItems, dependencies, report)
	s.detectSchemaInconsistencies(ctx, contracts, &report.Conflicts)
	s.detectContractCollisions(contracts, &report.Conflicts)
	s.detectVariableMismatches(variables, contracts, &report.Conflicts)
	s.checkValidationRulesConsistency(rules, report)

	// 4. Final scoring
	s.calculateAlignmentScore(report)

	// 5. Persist and return
	if err := s.repo.CreateReport(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *alignmentService) AnalyzeSnapshot(ctx context.Context, snapshot domain.ProjectSnapshot) (*domain.AlignmentReport, error) {
	report := &domain.AlignmentReport{
		ID:                     uuid.New(),
		Conflicts:              []domain.Conflict{},
		Overlaps:               []domain.Overlap{},
		MissingDependencies:    []string{},
		CircularDependencies:   []string{},
		RecommendedResolutions: []string{},
		AlignmentScore:         100,
		CreatedAt:              time.Now(),
	}

	// 1. Detect circular dependencies in Architecture
	s.detectSnapshotDependencyLoops(snapshot.Architecture, report)

	// 2. Detect schema inconsistencies in Contracts
	s.detectSnapshotSchemaInconsistencies(snapshot.Contracts, &report.Conflicts)

	// 3. Score
	s.calculateAlignmentScore(report)

	return report, nil
}

func (s *alignmentService) detectSnapshotDependencyLoops(arch domain.ArchitectureModel, report *domain.AlignmentReport) {
	adj := make(map[string][]string)
	for _, edge := range arch.DependenciesGraph {
		adj[edge.From] = append(adj[edge.From], edge.To)
	}

	visited := make(map[string]bool)
	onStack := make(map[string]bool)

	var hasCycle func(u string) bool
	hasCycle = func(u string) bool {
		visited[u] = true
		onStack[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				if hasCycle(v) {
					return true
				}
			} else if onStack[v] {
				return true
			}
		}
		onStack[u] = false
		return false
	}

	for _, layer := range arch.Layers {
		if !visited[layer] {
			if hasCycle(layer) {
				report.CircularDependencies = append(report.CircularDependencies, "Circular dependency detected in architecture layers involving: "+layer)
				report.Conflicts = append(report.Conflicts, domain.Conflict{
					ID:          uuid.New(),
					Severity:    domain.SeverityCritical,
					Type:        domain.ConflictDependencyLoop,
					Description: "Architectural cycle detected at layer: " + layer,
					Remediation: "Review layer isolation rules and break the coupling.",
					CreatedAt:   time.Now(),
				})
			}
		}
	}
}

func (s *alignmentService) detectSnapshotSchemaInconsistencies(contracts domain.ContractsModel, conflicts *[]domain.Conflict) {
	// Check for field name/type collisions across all API contracts and data models
	fieldTypes := make(map[string]string)
	fieldSources := make(map[string]string)

	for _, c := range contracts.APIContracts {
		// OutputSchema check (simplified)
		if c.ResponseSchema != nil {
			if props, ok := c.ResponseSchema["properties"].(map[string]interface{}); ok {
				for fieldName, fieldData := range props {
					if fd, ok := fieldData.(map[string]interface{}); ok {
						fType, _ := fd["type"].(string)
						if existingType, exists := fieldTypes[fieldName]; exists && existingType != fType {
							*conflicts = append(*conflicts, domain.Conflict{
								ID:          uuid.New(),
								Severity:    domain.SeverityError,
								Type:        domain.ConflictSchemaMismatch,
								Description: fmt.Sprintf("Field '%s' (type %s) in contract '%s' conflicts with type %s in '%s'.", fieldName, fType, c.Name, existingType, fieldSources[fieldName]),
								Remediation: "Normalize field types across the API surface.",
								CreatedAt:   time.Now(),
							})
						} else {
							fieldTypes[fieldName] = fType
							fieldSources[fieldName] = c.Name
						}
					}
				}
			}
		}
	}
}

func (s *alignmentService) detectSchemaInconsistencies(ctx context.Context, contracts []domain.ContractDefinition, conflicts *[]domain.Conflict) {
	// Example: Two contracts define the same field name in different types (very simplified)
	fieldTypes := make(map[string]string)
	fieldSources := make(map[string]uuid.UUID)

	for _, c := range contracts {
		// Traverse schemas (simplified check for top-level fields in output_schema)
		if c.OutputSchema != nil {
			if props, ok := c.OutputSchema["properties"].(map[string]interface{}); ok {
				for fieldName, fieldData := range props {
					if fd, ok := fieldData.(map[string]interface{}); ok {
						fType, _ := fd["type"].(string)
						if existingType, exists := fieldTypes[fieldName]; exists && existingType != fType {
							*conflicts = append(*conflicts, domain.Conflict{
								ID:          uuid.New(),
								Severity:    domain.SeverityError,
								Type:        domain.ConflictSchemaMismatch,
								SourceID:    fieldSources[fieldName],
								TargetID:    c.ID,
								Description: fmt.Sprintf("Field '%s' defined as %s in one contract and %s in another.", fieldName, existingType, fType),
								Remediation: fmt.Sprintf("Align the type of '%s' to be consistent across all contracts.", fieldName),
								CreatedAt:   time.Now(),
							})
						} else {
							fieldTypes[fieldName] = fType
							fieldSources[fieldName] = c.ID
						}
					}
				}
			}
		}
	}
}

func (s *alignmentService) detectDependencyLoops(items []domain.RoadmapItem, deps []domain.RoadmapDependency, report *domain.AlignmentReport) {
	// Build graph
	adj := make(map[uuid.UUID][]uuid.UUID)
	for _, d := range deps {
		adj[d.SourceID] = append(adj[d.SourceID], d.TargetID)
	}

	// DFS for cycle detection
	visited := make(map[uuid.UUID]bool)
	onStack := make(map[uuid.UUID]bool)

	var hasCycle func(u uuid.UUID) bool
	hasCycle = func(u uuid.UUID) bool {
		visited[u] = true
		onStack[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				if hasCycle(v) {
					return true
				}
			} else if onStack[v] {
				return true
			}
		}
		onStack[u] = false
		return false
	}

	for _, item := range items {
		if !visited[item.ID] {
			if hasCycle(item.ID) {
				report.CircularDependencies = append(report.CircularDependencies, "Circular dependency detected starting from: "+item.Title)
				report.Conflicts = append(report.Conflicts, domain.Conflict{
					ID:          uuid.New(),
					Severity:    domain.SeverityCritical,
					Type:        domain.ConflictDependencyLoop,
					SourceID:    item.ID,
					Description: "Circular dependency detected involving roadmap item: " + item.Title,
					Remediation: "Break the dependency loop by removing redundant dependencies.",
					CreatedAt:   time.Now(),
				})
			}
		}
	}
}

func (s *alignmentService) detectContractCollisions(contracts []domain.ContractDefinition, conflicts *[]domain.Conflict) {
	// Collision: Multiple contracts of the same type for the same Roadmap Item (unless versioned/documented differently)
	contractKeys := make(map[string]uuid.UUID)

	for _, c := range contracts {
		key := fmt.Sprintf("%s-%s", c.RoadmapItemID, c.ContractType)
		if originalID, exists := contractKeys[key]; exists && originalID != c.ID {
			*conflicts = append(*conflicts, domain.Conflict{
				ID:          uuid.New(),
				Severity:    domain.SeverityWarning,
				Type:        domain.ConflictContractCollision,
				SourceID:    originalID,
				TargetID:    c.ID,
				Description: fmt.Sprintf("Multiple %s contracts detected for the same Roadmap Item. This may cause integration ambiguity.", c.ContractType),
				Remediation: "Consolidate multiple contracts or clearly define versioning pathways.",
				CreatedAt:   time.Now(),
			})
		}
		contractKeys[key] = c.ID
	}
}

func (s *alignmentService) detectVariableMismatches(variables []domain.VariableDefinition, contracts []domain.ContractDefinition, conflicts *[]domain.Conflict) {
	// Check if variables defined in contracts actually exist in the variable registry
	// (Simplified check: for now just ensure they aren't orphaned)
	contractIDs := make(map[uuid.UUID]bool)
	for _, c := range contracts {
		contractIDs[c.ID] = true
	}

	for _, v := range variables {
		if !contractIDs[v.ContractID] {
			*conflicts = append(*conflicts, domain.Conflict{
				ID:          uuid.New(),
				Severity:    domain.SeverityWarning,
				Type:        domain.ConflictLogicContradiction,
				SourceID:    v.ID,
				Description: fmt.Sprintf("Variable '%s' is associated with a non-existent or deleted contract.", v.Name),
				Remediation: "Update the variable's contract association or remove the orphaned variable.",
				CreatedAt:   time.Now(),
			})
		}
	}
}

func (s *alignmentService) checkValidationRulesConsistency(rules []domain.ValidationRule, report *domain.AlignmentReport) {
	// Basic check: Ensure no duplicate rules by name
	ruleNames := make(map[string]bool)
	for _, r := range rules {
		if ruleNames[r.Name] {
			report.Conflicts = append(report.Conflicts, domain.Conflict{
				ID:          uuid.New(),
				Severity:    domain.SeverityInfo,
				Type:        domain.ConflictLogicContradiction,
				SourceID:    r.ID,
				Description: fmt.Sprintf("Duplicate validation rule name detected: %s", r.Name),
				Remediation: "Rename or consolidate duplicate validation rules.",
				CreatedAt:   time.Now(),
			})
		}
		ruleNames[r.Name] = true
	}
}

func (s *alignmentService) calculateAlignmentScore(report *domain.AlignmentReport) {
	score := 100
	for _, c := range report.Conflicts {
		switch c.Severity {
		case domain.SeverityCritical:
			score -= 30
		case domain.SeverityError:
			score -= 15
		case domain.SeverityWarning:
			score -= 5
		}
	}
	if score < 0 {
		score = 0
	}
	report.AlignmentScore = score
}
