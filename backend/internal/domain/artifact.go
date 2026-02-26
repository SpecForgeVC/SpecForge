package domain

import (
	"time"

	"github.com/google/uuid"
)

type ExportFormat string

const (
	ExportFormatJSON     ExportFormat = "json"
	ExportFormatMarkdown ExportFormat = "markdown"
	ExportFormatZip      ExportFormat = "zip"
)

type BuildArtifactPackage struct {
	Metadata              MetadataSection      `json:"metadata"`
	RoadmapContext        RoadmapContext       `json:"roadmapContext"`
	Contracts             []ContractBundle     `json:"contracts"`
	Schemas               []SchemaBundle       `json:"schemas"`
	ValidationRules       []ValidationBundle   `json:"validationRules"`
	Variables             []VariableBundle     `json:"variables"`
	Dependencies          DependencyGraph      `json:"dependencies"`
	GovernanceConstraints GovernanceBundle     `json:"governanceConstraints"`
	AlignmentReport       *AlignmentReport     `json:"alignmentReport,omitempty"`
	AcceptanceCriteria    []AcceptanceCriteria `json:"acceptanceCriteria"`
	TestRequirements      []TestSpecification  `json:"testRequirements"`
	BuildPrompts          BuildPromptBundle    `json:"buildPrompts"`
	RefinementLoopPrompts RefinementLoopBundle `json:"refinementLoopPrompts"`
}

type MetadataSection struct {
	ArtifactID     uuid.UUID `json:"artifact_id"`
	RoadmapItemID  uuid.UUID `json:"roadmap_item_id"`
	Version        string    `json:"version"`
	ExportedAt     time.Time `json:"exported_at"`
	ExportedBy     uuid.UUID `json:"exported_by"`
	IntegrityHash  string    `json:"integrity_hash"`
	GovernanceMode string    `json:"governance_mode"`
}

type RoadmapContext struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	BusinessContext  string `json:"business_context"`
	TechnicalContext string `json:"technical_context"`
	Priority         string `json:"priority"`
	RiskLevel        string `json:"risk_level"`
}

type ContractBundle struct {
	ID           uuid.UUID              `json:"id"`
	Type         string                 `json:"type"`
	Version      string                 `json:"version"`
	InputSchema  map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema"`
}

type SchemaBundle struct {
	Name       string                 `json:"name"`
	Definition map[string]interface{} `json:"definition"`
}

type ValidationBundle struct {
	Name     string                 `json:"name"`
	RuleType string                 `json:"rule_type"`
	Config   map[string]interface{} `json:"config"`
}

type VariableBundle struct {
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Required        bool                   `json:"required"`
	ValidationRules map[string]interface{} `json:"validation_rules"`
}

type DependencyGraph struct {
	Nodes []string         `json:"nodes"`
	Edges []DependencyEdge `json:"edges"`
}

type DependencyEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"` // enum: contract, schema, validation, roadmap
}

type GovernanceBundle struct {
	PoliciesEnforced []string `json:"policies_enforced"`
	ComplianceStatus string   `json:"compliance_status"`
}

type AcceptanceCriteria struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type TestSpecification struct {
	Type        string `json:"type"`
	Instruction string `json:"instruction"`
}

type BuildPromptBundle struct {
	Implementation string `json:"implementation"`
	Verification   string `json:"verification"`
}

type RefinementLoopBundle struct {
	Instructions string `json:"instructions"`
}
