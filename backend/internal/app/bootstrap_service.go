package app

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type bootstrapService struct {
	repo        BootstrapRepository
	projectRepo ProjectRepository
	sessionRepo ImportSessionRepository
}

func NewBootstrapService(repo BootstrapRepository, projectRepo ProjectRepository, sessionRepo ImportSessionRepository) BootstrapService {
	return &bootstrapService{
		repo:        repo,
		projectRepo: projectRepo,
		sessionRepo: sessionRepo,
	}
}

const ideAnalysisPromptTemplate = `You are analyzing an existing codebase to create a structured intelligence snapshot for SpecForge.

**Project: %s**

Analyze the entire codebase and produce a JSON object with the following sections. Be thorough and precise.

## Output Format (strict JSON)

{
  "project_overview": {
    "name": "string",
    "description": "string",
    "domain": "string",
    "primary_language": "string",
    "architecture_pattern": "string (e.g. MVC, Clean Architecture, Microservices)"
  },
  "tech_stack": {
    "languages": ["string"],
    "frameworks": ["string"],
    "databases": ["string"],
    "infrastructure": ["string"],
    "build_tools": ["string"]
  },
  "modules": [
    {
      "name": "string",
      "description": "string",
      "responsibilities": ["string"],
      "risk_level": "LOW | MEDIUM | HIGH",
      "change_sensitivity": "LOW | MEDIUM | HIGH",
      "dependencies": ["module_name"]
    }
  ],
  "apis": [
    {
      "endpoint": "string",
      "method": "GET | POST | PUT | PATCH | DELETE",
      "auth_type": "string (e.g. JWT, API_KEY, NONE)",
      "request_schema": {},
      "response_schema": {}
    }
  ],
  "data_models": [
    {
      "name": "string",
      "relationships": [
        {"target": "string", "type": "string (ONE_TO_MANY, MANY_TO_ONE, etc.)"}
      ],
      "constraints": [
        {"field": "string", "rule": "string"}
      ]
    }
  ],
  "validation_rules": [
    {
      "name": "string",
      "target": "string",
      "rule_type": "string",
      "config": {}
    }
  ],
  "contracts": [
    {
      "name": "string",
      "contract_type": "REST | GRAPHQL | EVENT | INTERNAL_FUNCTION",
      "schema": {},
      "source_module": "string",
      "stability_score": 0.0
    }
  ],
  "current_state": {
    "test_coverage_estimate": "string",
    "documentation_level": "string",
    "code_quality_notes": "string",
    "known_tech_debt": ["string"]
  },
  "risks": [
    {
      "area": "string",
      "severity": "LOW | MEDIUM | HIGH | CRITICAL",
      "description": "string"
    }
  ],
  "change_sensitivity": [
    {
      "module": "string",
      "sensitivity": "LOW | MEDIUM | HIGH",
      "reason": "string"
    }
  ]
}

IMPORTANT:
- Output ONLY the JSON object, no markdown, no explanation
- Be exhaustive — include ALL endpoints, models, and modules
- For risk_level and change_sensitivity, consider coupling, complexity, and blast radius
- Stability scores should be between 0.0 (unstable) and 1.0 (stable)
`

func (s *bootstrapService) GeneratePrompt(ctx context.Context, projectID uuid.UUID) (string, string, error) {
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return "", "", fmt.Errorf("get project: %w", err)
	}
	prompt := fmt.Sprintf(ideAnalysisPromptTemplate, project.Name)
	return prompt, project.Name, nil
}

func (s *bootstrapService) IngestBootstrap(ctx context.Context, projectID uuid.UUID, payload domain.BootstrapPayload) (*domain.ProjectIntelligenceSnapshot, []string, error) {
	var warnings []string

	// Validate minimum structure
	if payload.ProjectOverview == nil {
		warnings = append(warnings, "project_overview section is empty")
	}
	if len(payload.Modules) == 0 {
		warnings = append(warnings, "no modules detected in analysis")
	}
	if len(payload.APIs) == 0 {
		warnings = append(warnings, "no APIs detected in analysis")
	}

	// Get next version
	maxVersion, err := s.repo.GetMaxVersion(ctx, projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("get max version: %w", err)
	}

	// Build snapshot JSON from payload
	payloadBytes, _ := json.Marshal(payload)
	var snapshotJSON map[string]interface{}
	json.Unmarshal(payloadBytes, &snapshotJSON)

	// Compute scores
	scores := s.computeScores(payload)

	// Build confidence map
	confidence := s.computeConfidence(payload)

	// Create snapshot
	snapshot := &domain.ProjectIntelligenceSnapshot{
		ProjectID:         projectID,
		Version:           maxVersion + 1,
		SnapshotJSON:      snapshotJSON,
		ArchitectureScore: scores.ArchitectureScore,
		ContractDensity:   scores.ContractDensity,
		RiskScore:         scores.RiskScore,
		AlignmentScore:    scores.AlignmentScore,
		ConfidenceJSON:    confidence,
	}

	if err := s.repo.InsertSnapshot(ctx, snapshot); err != nil {
		return nil, nil, fmt.Errorf("insert snapshot: %w", err)
	}

	// Insert normalized sub-resources
	// Modules
	for _, m := range payload.Modules {
		module := &domain.ProjectModule{
			ProjectID:         projectID,
			SnapshotID:        snapshot.ID,
			Name:              getStr(m, "name"),
			Description:       getStr(m, "description"),
			RiskLevel:         strings.ToUpper(getStr(m, "risk_level")),
			ChangeSensitivity: strings.ToUpper(getStr(m, "change_sensitivity")),
		}
		if err := s.repo.InsertModule(ctx, module); err != nil {
			warnings = append(warnings, fmt.Sprintf("failed to insert module %s: %v", module.Name, err))
		}
	}

	// Data models → Entities
	for _, dm := range payload.DataModels {
		entity := &domain.ProjectEntity{
			ProjectID:  projectID,
			SnapshotID: snapshot.ID,
			Name:       getStr(dm, "name"),
		}
		if rels, ok := dm["relationships"]; ok {
			relBytes, _ := json.Marshal(rels)
			var relMap map[string]interface{}
			json.Unmarshal(relBytes, &relMap)
			entity.RelationshipsJSON = relMap
		}
		if cons, ok := dm["constraints"]; ok {
			conBytes, _ := json.Marshal(cons)
			var conMap map[string]interface{}
			json.Unmarshal(conBytes, &conMap)
			entity.ConstraintsJSON = conMap
		}
		if err := s.repo.InsertEntity(ctx, entity); err != nil {
			warnings = append(warnings, fmt.Sprintf("failed to insert entity %s: %v", entity.Name, err))
		}
	}

	// APIs
	for _, a := range payload.APIs {
		entry := &domain.ProjectApiEntry{
			ProjectID:  projectID,
			SnapshotID: snapshot.ID,
			Endpoint:   getStr(a, "endpoint"),
			Method:     strings.ToUpper(getStr(a, "method")),
			AuthType:   getStr(a, "auth_type"),
		}
		if rs, ok := a["request_schema"]; ok {
			rsBytes, _ := json.Marshal(rs)
			var rsMap map[string]interface{}
			json.Unmarshal(rsBytes, &rsMap)
			entry.RequestSchema = rsMap
		}
		if rs, ok := a["response_schema"]; ok {
			rsBytes, _ := json.Marshal(rs)
			var rsMap map[string]interface{}
			json.Unmarshal(rsBytes, &rsMap)
			entry.ResponseSchema = rsMap
		}
		if err := s.repo.InsertApiEntry(ctx, entry); err != nil {
			warnings = append(warnings, fmt.Sprintf("failed to insert api %s: %v", entry.Endpoint, err))
		}
	}

	// Contracts
	for _, c := range payload.Contracts {
		contractEntry := &domain.ProjectContractEntry{
			ProjectID:    projectID,
			SnapshotID:   snapshot.ID,
			Name:         getStr(c, "name"),
			ContractType: getStr(c, "contract_type"),
			SourceModule: getStr(c, "source_module"),
		}
		if schema, ok := c["schema"]; ok {
			schemaBytes, _ := json.Marshal(schema)
			var schemaMap map[string]interface{}
			json.Unmarshal(schemaBytes, &schemaMap)
			contractEntry.SchemaJSON = schemaMap
		}
		if score, ok := c["stability_score"]; ok {
			if f, ok := score.(float64); ok {
				contractEntry.StabilityScore = f
			}
		}
		if err := s.repo.InsertContractEntry(ctx, contractEntry); err != nil {
			warnings = append(warnings, fmt.Sprintf("failed to insert contract %s: %v", contractEntry.Name, err))
		}
	}

	return snapshot, warnings, nil
}

func (s *bootstrapService) ListSnapshots(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectIntelligenceSnapshot, error) {
	return s.repo.ListSnapshots(ctx, projectID)
}

func (s *bootstrapService) GetSnapshot(ctx context.Context, snapshotID uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error) {
	return s.repo.GetSnapshot(ctx, snapshotID)
}

func (s *bootstrapService) GetLatestSnapshot(ctx context.Context, projectID uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error) {
	return s.repo.GetLatestSnapshot(ctx, projectID)
}

func (s *bootstrapService) GetLatestImportSession(ctx context.Context, projectID uuid.UUID) (*domain.ImportSession, error) {
	return s.sessionRepo.GetLatestSessionByProject(ctx, projectID)
}

func (s *bootstrapService) DiffSnapshots(ctx context.Context, projectID uuid.UUID, fromID, toID *uuid.UUID) (map[string]interface{}, error) {
	var fromSnap, toSnap *domain.ProjectIntelligenceSnapshot
	var err error

	if fromID != nil {
		fromSnap, err = s.repo.GetSnapshot(ctx, *fromID)
		if err != nil {
			return nil, fmt.Errorf("get from snapshot: %w", err)
		}
	}

	if toID != nil {
		toSnap, err = s.repo.GetSnapshot(ctx, *toID)
		if err != nil {
			return nil, fmt.Errorf("get to snapshot: %w", err)
		}
	} else {
		toSnap, err = s.repo.GetLatestSnapshot(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("get latest snapshot: %w", err)
		}
	}

	// Build diff result
	diff := map[string]interface{}{
		"from_snapshot_id": nil,
		"to_snapshot_id":   toSnap.ID,
		"from_version":     0,
		"to_version":       toSnap.Version,
	}

	if fromSnap != nil {
		diff["from_snapshot_id"] = fromSnap.ID
		diff["from_version"] = fromSnap.Version

		// Compare modules
		fromModules, _ := s.repo.ListModulesBySnapshot(ctx, fromSnap.ID)
		toModules, _ := s.repo.ListModulesBySnapshot(ctx, toSnap.ID)
		diff["added_modules"], diff["removed_modules"] = diffByName(fromModules, toModules)

		// Compare APIs
		fromApis, _ := s.repo.ListApiEntriesBySnapshot(ctx, fromSnap.ID)
		toApis, _ := s.repo.ListApiEntriesBySnapshot(ctx, toSnap.ID)
		diff["added_apis"], diff["removed_apis"] = diffApisByEndpoint(fromApis, toApis)

		// Compare entities
		fromEntities, _ := s.repo.ListEntitiesBySnapshot(ctx, fromSnap.ID)
		toEntities, _ := s.repo.ListEntitiesBySnapshot(ctx, toSnap.ID)
		diff["added_entities"], diff["removed_entities"] = diffEntitiesByName(fromEntities, toEntities)

		// Score changes
		diff["score_changes"] = map[string]interface{}{
			"architecture_score": toSnap.ArchitectureScore - fromSnap.ArchitectureScore,
			"contract_density":   toSnap.ContractDensity - fromSnap.ContractDensity,
			"risk_score":         toSnap.RiskScore - fromSnap.RiskScore,
			"alignment_score":    toSnap.AlignmentScore - fromSnap.AlignmentScore,
		}
	}

	return diff, nil
}

// --- Score computation ---

func (s *bootstrapService) computeScores(payload domain.BootstrapPayload) domain.BootstrapScores {
	moduleCount := float64(len(payload.Modules))
	apiCount := float64(len(payload.APIs))
	contractCount := float64(len(payload.Contracts))
	entityCount := float64(len(payload.DataModels))

	// Architecture complexity: based on module count and inter-dependencies
	archScore := math.Min(100, moduleCount*10+apiCount*2)

	// Contract density: contracts per module
	contractDensity := float64(0)
	if moduleCount > 0 {
		contractDensity = math.Min(100, (contractCount/moduleCount)*50)
	}

	// Risk score: based on high-risk items
	highRiskCount := float64(0)
	for _, r := range payload.Risks {
		if sev, ok := r["severity"]; ok {
			if strings.ToUpper(fmt.Sprintf("%v", sev)) == "HIGH" || strings.ToUpper(fmt.Sprintf("%v", sev)) == "CRITICAL" {
				highRiskCount++
			}
		}
	}
	riskScore := math.Min(100, highRiskCount*20)

	// Alignment: based on completeness of analysis
	sectionsFilled := float64(0)
	if payload.ProjectOverview != nil {
		sectionsFilled++
	}
	if payload.TechStack != nil {
		sectionsFilled++
	}
	if len(payload.Modules) > 0 {
		sectionsFilled++
	}
	if len(payload.APIs) > 0 {
		sectionsFilled++
	}
	if len(payload.DataModels) > 0 {
		sectionsFilled++
	}
	if len(payload.Contracts) > 0 {
		sectionsFilled++
	}
	if len(payload.Risks) > 0 {
		sectionsFilled++
	}
	alignmentScore := (sectionsFilled / 7.0) * 100

	_ = entityCount // used for potential future scoring

	return domain.BootstrapScores{
		ArchitectureScore: math.Round(archScore*100) / 100,
		ContractDensity:   math.Round(contractDensity*100) / 100,
		RiskScore:         math.Round(riskScore*100) / 100,
		AlignmentScore:    math.Round(alignmentScore*100) / 100,
	}
}

func (s *bootstrapService) computeConfidence(payload domain.BootstrapPayload) map[string]interface{} {
	conf := map[string]interface{}{}
	conf["project_overview"] = boolToConfidence(len(payload.ProjectOverview) > 0)
	conf["tech_stack"] = boolToConfidence(len(payload.TechStack) > 0)
	conf["modules"] = arrayConfidence(len(payload.Modules))
	conf["apis"] = arrayConfidence(len(payload.APIs))
	conf["data_models"] = arrayConfidence(len(payload.DataModels))
	conf["validation_rules"] = arrayConfidence(len(payload.ValidationRules))
	conf["contracts"] = arrayConfidence(len(payload.Contracts))
	conf["current_state"] = boolToConfidence(len(payload.CurrentState) > 0)
	conf["risks"] = arrayConfidence(len(payload.Risks))
	conf["change_sensitivity"] = arrayConfidence(len(payload.ChangeSensivity))
	return conf
}

func boolToConfidence(present bool) float64 {
	if present {
		return 1.0
	}
	return 0.0
}

func arrayConfidence(count int) float64 {
	if count == 0 {
		return 0.0
	}
	if count >= 5 {
		return 1.0
	}
	return float64(count) / 5.0
}

// --- Helpers ---

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

type named interface {
	GetName() string
}

func diffByName(from, to []domain.ProjectModule) ([]domain.ProjectModule, []domain.ProjectModule) {
	fromSet := make(map[string]bool)
	toSet := make(map[string]bool)
	for _, m := range from {
		fromSet[m.Name] = true
	}
	for _, m := range to {
		toSet[m.Name] = true
	}

	var added, removed []domain.ProjectModule
	for _, m := range to {
		if !fromSet[m.Name] {
			added = append(added, m)
		}
	}
	for _, m := range from {
		if !toSet[m.Name] {
			removed = append(removed, m)
		}
	}
	return added, removed
}

func diffApisByEndpoint(from, to []domain.ProjectApiEntry) ([]domain.ProjectApiEntry, []domain.ProjectApiEntry) {
	key := func(a domain.ProjectApiEntry) string { return a.Method + " " + a.Endpoint }
	fromSet := make(map[string]bool)
	toSet := make(map[string]bool)
	for _, a := range from {
		fromSet[key(a)] = true
	}
	for _, a := range to {
		toSet[key(a)] = true
	}

	var added, removed []domain.ProjectApiEntry
	for _, a := range to {
		if !fromSet[key(a)] {
			added = append(added, a)
		}
	}
	for _, a := range from {
		if !toSet[key(a)] {
			removed = append(removed, a)
		}
	}
	return added, removed
}

func diffEntitiesByName(from, to []domain.ProjectEntity) ([]domain.ProjectEntity, []domain.ProjectEntity) {
	fromSet := make(map[string]bool)
	toSet := make(map[string]bool)
	for _, e := range from {
		fromSet[e.Name] = true
	}
	for _, e := range to {
		toSet[e.Name] = true
	}

	var added, removed []domain.ProjectEntity
	for _, e := range to {
		if !fromSet[e.Name] {
			added = append(added, e)
		}
	}
	for _, e := range from {
		if !toSet[e.Name] {
			removed = append(removed, e)
		}
	}
	return added, removed
}
