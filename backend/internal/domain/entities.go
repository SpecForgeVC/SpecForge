package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner    Role = "OWNER"
	RoleAdmin    Role = "ADMIN"
	RoleReviewer Role = "REVIEWER"
	RoleEngineer Role = "ENGINEER"
	RoleAIAgent  Role = "AI_AGENT"
)

type Principal struct {
	UserID      uuid.UUID `json:"user_id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	Role        Role      `json:"role"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         Role      `json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Workspace struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MCPSettings struct {
	Enabled        bool       `json:"enabled"`
	Port           int        `json:"port"`
	BindAddress    string     `json:"bind_address"`
	TokenRequired  bool       `json:"token_required"`
	Token          string     `json:"token"`
	HealthStatus   string     `json:"health_status"`
	LastSnapshotAt *time.Time `json:"last_snapshot_at,omitempty"`
	ImportMode     string     `json:"import_mode"`     // light, full
	RepositoryType string     `json:"repository_type"` // monorepo, polyrepo, single
}

type Project struct {
	ID            uuid.UUID              `json:"id"`
	WorkspaceID   uuid.UUID              `json:"workspace_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	TechStack     map[string]interface{} `json:"tech_stack"`
	Settings      map[string]interface{} `json:"settings"`
	MCPSettings   MCPSettings            `json:"mcp_settings"`
	RepositoryURL string                 `json:"repository_url"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type TechStackRecommendation struct {
	RecommendedStack map[string]interface{} `json:"recommended_stack"`
	Reasoning        string                 `json:"reasoning"`
}

type RoadmapItemType string

const (
	Epic     RoadmapItemType = "EPIC"
	Feature  RoadmapItemType = "FEATURE"
	Task     RoadmapItemType = "TASK"
	Bugfix   RoadmapItemType = "BUGFIX"
	Refactor RoadmapItemType = "REFACTOR"
)

type RoadmapItemPriority string

const (
	PriorityLow      RoadmapItemPriority = "LOW"
	PriorityMedium   RoadmapItemPriority = "MEDIUM"
	PriorityHigh     RoadmapItemPriority = "HIGH"
	PriorityCritical RoadmapItemPriority = "CRITICAL"
)

type RoadmapItemStatus string

const (
	StatusDraft      RoadmapItemStatus = "DRAFT"
	StatusInReview   RoadmapItemStatus = "IN_REVIEW"
	StatusApproved   RoadmapItemStatus = "APPROVED"
	StatusInProgress RoadmapItemStatus = "IN_PROGRESS"
	StatusComplete   RoadmapItemStatus = "COMPLETE"
)

type APIRecommendation struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	TechnicalContext string `json:"technical_context"`
}

type RiskLevel string

const (
	RiskLow    RiskLevel = "LOW"
	RiskMedium RiskLevel = "MEDIUM"
	RiskHigh   RiskLevel = "HIGH"
)

type RoadmapItem struct {
	ID                  uuid.UUID           `json:"id"`
	ProjectID           uuid.UUID           `json:"project_id"`
	Type                RoadmapItemType     `json:"type"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	BusinessContext     string              `json:"business_context"`
	TechnicalContext    string              `json:"technical_context"`
	Priority            RoadmapItemPriority `json:"priority"`
	Status              RoadmapItemStatus   `json:"status"`
	RiskLevel           RiskLevel           `json:"risk_level"`
	ReadinessLevel      ReadinessLevel      `json:"readiness_level"`
	BreakingChange      bool                `json:"breaking_change"`
	RegressionSensitive bool                `json:"regression_sensitive"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

type ContractType string

const (
	REST             ContractType = "REST"
	GraphQL          ContractType = "GRAPHQL"
	CLI              ContractType = "CLI"
	InternalFunction ContractType = "INTERNAL_FUNCTION"
	Event            ContractType = "EVENT"
)

type ContractDefinition struct {
	ID                 uuid.UUID              `json:"id"`
	RoadmapItemID      uuid.UUID              `json:"roadmap_item_id"`
	ContractType       ContractType           `json:"contract_type"`
	Version            string                 `json:"version"`
	InputSchema        map[string]interface{} `json:"input_schema"`
	OutputSchema       map[string]interface{} `json:"output_schema"`
	ErrorSchema        map[string]interface{} `json:"error_schema"`
	BackwardCompatible bool                   `json:"backward_compatible"`
	DeprecatedFields   []string               `json:"deprecated_fields"`
	CreatedAt          time.Time              `json:"created_at"`
}

type VersionSnapshot struct {
	ID            uuid.UUID              `json:"id"`
	RoadmapItemID uuid.UUID              `json:"roadmap_item_id"`
	SnapshotData  map[string]interface{} `json:"snapshot_data"`
	Hash          string                 `json:"hash"`
	CreatedAt     time.Time              `json:"created_at"`
	CreatedBy     uuid.UUID              `json:"created_by"`
}

type ProposalType string

const (
	EditDescription ProposalType = "EDIT_DESCRIPTION"
	ModifySchema    ProposalType = "MODIFY_SCHEMA"
	AddVariable     ProposalType = "ADD_VARIABLE"
	RemoveField     ProposalType = "REMOVE_FIELD"
)

type ProposalStatus string

const (
	Pending  ProposalStatus = "PENDING"
	Approved ProposalStatus = "APPROVED"
	Rejected ProposalStatus = "REJECTED"
)

type AiProposal struct {
	ID              uuid.UUID              `json:"id"`
	RoadmapItemID   uuid.UUID              `json:"roadmap_item_id"`
	ProposalType    ProposalType           `json:"proposal_type"`
	Diff            map[string]interface{} `json:"diff"`
	Reasoning       string                 `json:"reasoning"`
	ConfidenceScore float64                `json:"confidence_score"`
	Status          ProposalStatus         `json:"status"`
	ReviewedBy      uuid.UUID              `json:"reviewed_by,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

type AuditLog struct {
	ID          uuid.UUID              `json:"id"`
	EntityType  string                 `json:"entity_type"`
	EntityID    uuid.UUID              `json:"entity_id"`
	Action      string                 `json:"action"`
	PerformedBy uuid.UUID              `json:"performed_by,omitempty"`
	OldData     map[string]interface{} `json:"old_data"`
	NewData     map[string]interface{} `json:"new_data"`
	CreatedAt   time.Time              `json:"created_at"`
}

type Requirement struct {
	ID                 uuid.UUID `json:"id"`
	RoadmapItemID      uuid.UUID `json:"roadmap_item_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Testable           bool      `json:"testable"`
	AcceptanceCriteria string    `json:"acceptance_criteria"`
	OrderIndex         int       `json:"order_index"`
}

type VariableDefinition struct {
	ID              uuid.UUID              `json:"id"`
	ContractID      uuid.UUID              `json:"contract_id"`
	Name            string                 `json:"name"`
	Type            string                 `json:"type"`
	Required        bool                   `json:"required"`
	DefaultValue    string                 `json:"default_value"`
	Description     string                 `json:"description"`
	ValidationRules map[string]interface{} `json:"validation_rules"`
}

type ValidationRule struct {
	ID          uuid.UUID              `json:"id"`
	ProjectID   uuid.UUID              `json:"project_id"`
	Name        string                 `json:"name"`
	RuleType    string                 `json:"rule_type"`
	RuleConfig  map[string]interface{} `json:"rule_config"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
}

type Webhook struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Secret    string    `json:"secret"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type DriftReport struct {
	DriftDetected   bool             `json:"drift_detected"`
	BreakingChanges []BreakingChange `json:"breaking_changes"`
	RiskScore       float64          `json:"risk_score"`
}

type BreakingChange struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

type DriftFix struct {
	Field           string `json:"field"`
	Issue           string `json:"issue"`
	SuggestedChange string `json:"suggested_change"`
	Explanation     string `json:"explanation"`
}
type FeatureIntelligence struct {
	ID                       uuid.UUID `json:"id"`
	FeatureID                uuid.UUID `json:"feature_id"`
	CompletenessScore        int       `json:"completeness_score"`
	ContractIntegrityScore   int       `json:"contract_integrity_score"`
	VariableCoverageScore    int       `json:"variable_coverage_score"`
	DependencyStabilityScore int       `json:"dependency_stability_score"`
	DriftRiskScore           int       `json:"drift_risk_score"`
	TestCoverageScore        int       `json:"test_coverage_score"`
	LLMConfidenceScore       int       `json:"llm_confidence_score"`
	OverallScore             int       `json:"overall_score"`
	LastCalculatedAt         time.Time `json:"last_calculated_at"`
}

type LineageEventType string

const (
	LineageDeclared         LineageEventType = "DECLARED"
	LineageMutated          LineageEventType = "MUTATED"
	LineageTypeChanged      LineageEventType = "TYPE_CHANGED"
	LineageMappedToContract LineageEventType = "MAPPED_TO_CONTRACT"
	LineagePassedToAPI      LineageEventType = "PASSED_TO_API"
	LineageUsedInTest       LineageEventType = "USED_IN_TEST"
	LineageRemoved          LineageEventType = "REMOVED"
)

type VariableLineageEvent struct {
	ID              uuid.UUID        `json:"id"`
	VariableID      uuid.UUID        `json:"variable_id"`
	EventType       LineageEventType `json:"event_type"`
	SourceComponent string           `json:"source_component"`
	Description     string           `json:"description"`
	PerformedBy     uuid.UUID        `json:"performed_by,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	Metadata        map[string]any   `json:"metadata"`
}

type DependencyType string

const (
	DependencyDirect   DependencyType = "DIRECT"
	DependencyDerived  DependencyType = "DERIVED"
	DependencyContract DependencyType = "CONTRACT"
)

type VariableDependency struct {
	ID               uuid.UUID      `json:"id"`
	SourceVariableID uuid.UUID      `json:"source_variable_id"`
	TargetVariableID uuid.UUID      `json:"target_variable_id"`
	DependencyType   DependencyType `json:"dependency_type"`
	CreatedAt        time.Time      `json:"created_at"`
}

type ReadinessLevel string

const (
	ReadinessReady           ReadinessLevel = "READY"
	ReadinessReview          ReadinessLevel = "REVIEW"
	ReadinessNeedsRefinement ReadinessLevel = "NEEDS_REFINEMENT"
	ReadinessBlocked         ReadinessLevel = "BLOCKED"
)

// --- Project Bootstrap Intelligence ---

type ProjectIntelligenceSnapshot struct {
	ID                uuid.UUID              `json:"id"`
	ProjectID         uuid.UUID              `json:"project_id"`
	Version           int                    `json:"version"`
	SnapshotJSON      map[string]interface{} `json:"snapshot_json"`
	ArchitectureScore float64                `json:"architecture_score"`
	ContractDensity   float64                `json:"contract_density"`
	RiskScore         float64                `json:"risk_score"`
	AlignmentScore    float64                `json:"alignment_score"`
	ConfidenceJSON    map[string]interface{} `json:"confidence"`
	CreatedAt         time.Time              `json:"created_at"`
}

type ProjectModule struct {
	ID                uuid.UUID `json:"id"`
	ProjectID         uuid.UUID `json:"project_id"`
	SnapshotID        uuid.UUID `json:"snapshot_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	RiskLevel         string    `json:"risk_level"`
	ChangeSensitivity string    `json:"change_sensitivity"`
}

type ProjectEntity struct {
	ID                uuid.UUID              `json:"id"`
	ProjectID         uuid.UUID              `json:"project_id"`
	SnapshotID        uuid.UUID              `json:"snapshot_id"`
	Name              string                 `json:"name"`
	RelationshipsJSON map[string]interface{} `json:"relationships"`
	ConstraintsJSON   map[string]interface{} `json:"constraints"`
}

type ProjectApiEntry struct {
	ID             uuid.UUID              `json:"id"`
	ProjectID      uuid.UUID              `json:"project_id"`
	SnapshotID     uuid.UUID              `json:"snapshot_id"`
	Endpoint       string                 `json:"endpoint"`
	Method         string                 `json:"method"`
	AuthType       string                 `json:"auth_type"`
	RequestSchema  map[string]interface{} `json:"request_schema"`
	ResponseSchema map[string]interface{} `json:"response_schema"`
}

type ProjectContractEntry struct {
	ID             uuid.UUID              `json:"id"`
	ProjectID      uuid.UUID              `json:"project_id"`
	SnapshotID     uuid.UUID              `json:"snapshot_id"`
	Name           string                 `json:"name"`
	ContractType   string                 `json:"contract_type"`
	SchemaJSON     map[string]interface{} `json:"schema"`
	SourceModule   string                 `json:"source_module"`
	StabilityScore float64                `json:"stability_score"`
}

type BootstrapPayload struct {
	ProjectOverview map[string]interface{}   `json:"project_overview"`
	TechStack       map[string]interface{}   `json:"tech_stack"`
	Modules         []map[string]interface{} `json:"modules"`
	APIs            []map[string]interface{} `json:"apis"`
	DataModels      []map[string]interface{} `json:"data_models"`
	ValidationRules []map[string]interface{} `json:"validation_rules"`
	Contracts       []map[string]interface{} `json:"contracts"`
	CurrentState    map[string]interface{}   `json:"current_state"`
	Risks           []map[string]interface{} `json:"risks"`
	ChangeSensivity []map[string]interface{} `json:"change_sensitivity"`
}

type BootstrapScores struct {
	ArchitectureScore float64 `json:"architecture_score"`
	ContractDensity   float64 `json:"contract_density"`
	RiskScore         float64 `json:"risk_score"`
	AlignmentScore    float64 `json:"alignment_score"`
}

// --- Intelligence Alignment ---

type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityError    Severity = "ERROR"
	SeverityCritical Severity = "CRITICAL"
)

type ConflictType string

const (
	ConflictSchemaMismatch     ConflictType = "SCHEMA_MISMATCH"
	ConflictContractCollision  ConflictType = "CONTRACT_COLLISION"
	ConflictLogicContradiction ConflictType = "LOGIC_CONTRADICTION"
	ConflictDependencyLoop     ConflictType = "DEPENDENCY_LOOP"
)

type Conflict struct {
	ID          uuid.UUID    `json:"id"`
	Severity    Severity     `json:"severity"`
	Type        ConflictType `json:"type"`
	SourceID    uuid.UUID    `json:"source_id"`
	TargetID    uuid.UUID    `json:"target_id"`
	Description string       `json:"description"`
	Remediation string       `json:"remediation"`
	CreatedAt   time.Time    `json:"created_at"`
}

type Overlap struct {
	Type         string   `json:"type"`
	SharedFields []string `json:"shared_fields"`
	Description  string   `json:"description"`
}

type AlignmentReport struct {
	ID                     uuid.UUID  `json:"id"`
	ProjectID              uuid.UUID  `json:"project_id"`
	Conflicts              []Conflict `json:"conflicts"`
	Overlaps               []Overlap  `json:"overlaps"`
	MissingDependencies    []string   `json:"missing_dependencies"`
	CircularDependencies   []string   `json:"circular_dependencies"`
	RecommendedResolutions []string   `json:"recommended_resolutions"`
	AlignmentScore         int        `json:"alignment_score"`
	CreatedAt              time.Time  `json:"created_at"`
}

type RoadmapDependency struct {
	ID             uuid.UUID      `json:"id"`
	SourceID       uuid.UUID      `json:"source_id"`
	TargetID       uuid.UUID      `json:"target_id"`
	DependencyType DependencyType `json:"dependency_type"`
	CreatedAt      time.Time      `json:"created_at"`
}

// --- Project Import ---

type ProjectSnapshot struct {
	Metadata     ProjectMetadata   `json:"metadata"`
	Architecture ArchitectureModel `json:"architecture"`
	Contracts    ContractsModel    `json:"contracts"`
	Variables    VariablesModel    `json:"variables"`
	Security     SecurityModel     `json:"security"`
	Tests        TestModel         `json:"tests"`
}

type ProjectMetadata struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Languages      []string `json:"languages"`
	Frameworks     []string `json:"frameworks"`
	Entrypoints    []string `json:"entrypoints"`
	RuntimeType    string   `json:"runtime_type"`
	DeploymentType string   `json:"deployment_type"`
}

type ArchitectureModel struct {
	Layers            []string         `json:"layers"`
	ServiceBoundaries []string         `json:"service_boundaries"`
	DependenciesGraph []DependencyEdge `json:"dependencies_graph"`
}

type ContractsModel struct {
	APIContracts    []APIContract `json:"api_contracts"`
	DataModels      []DataModel   `json:"data_models"`
	ValidationRules []string      `json:"validation_rules"`
}

type APIContract struct {
	Name           string                 `json:"name"`
	Method         string                 `json:"method"`
	Path           string                 `json:"path"`
	RequestSchema  map[string]interface{} `json:"request_schema"`
	ResponseSchema map[string]interface{} `json:"response_schema"`
}

type DataModel struct {
	Name   string      `json:"name"`
	Fields []DataField `json:"fields"`
}

type DataField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type VariablesModel struct {
	GlobalConfig         []string `json:"global_config"`
	EnvironmentVariables []string `json:"environment_variables"`
	FeatureFlags         []string `json:"feature_flags"`
}

type SecurityModel struct {
	AuthMechanisms       []string `json:"auth_mechanisms"`
	Roles                []string `json:"roles"`
	ExternalIntegrations []string `json:"external_integrations"`
}

type TestModel struct {
	CoverageEstimate float64  `json:"coverage_estimate"`
	TestTypes        []string `json:"test_types"`
}

// --- MCP core types ---

type SnapshotState string

const (
	StateInitiated    SnapshotState = "initiated"
	StateAwaitingPost SnapshotState = "awaiting_post"
	StateAnalyzing    SnapshotState = "analyzing"
	StateCompleted    SnapshotState = "completed"
	StateFailed       SnapshotState = "failed"
)

type EnvironmentSnapshot struct {
	Metadata        SnapshotMetadata `json:"metadata"`
	FileTree        []string         `json:"file_tree"`
	Dependencies    Dependencies     `json:"dependencies"`
	Environment     Environment      `json:"environment"`
	APIRoutes       []APIRoute       `json:"api_routes"`
	Database        DatabaseInfo     `json:"database"`
	MiddlewareStack []string         `json:"middleware_stack"`
	TestFramework   TestFramework    `json:"test_framework"`
	CIConfig        CIConfig         `json:"ci_config"`
	ExportsSurface  []string         `json:"exports_surface"`
}

type SnapshotMetadata struct {
	Timestamp       time.Time `json:"timestamp"`
	IDEType         string    `json:"ide_type"`
	Language        string    `json:"language"`
	ConfidenceScore float64   `json:"confidence_score"`
}

type Dependencies struct {
	Runtime  []string          `json:"runtime"`
	Dev      []string          `json:"dev"`
	Versions map[string]string `json:"versions"`
}

type Environment struct {
	EnvVariables []string `json:"env_variables"`
	ConfigFiles  []string `json:"config_files"`
}

type APIRoute struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
	Params  []string `json:"params"`
}

type DatabaseInfo struct {
	Migrations     []string `json:"migrations"`
	DetectedTables []string `json:"detected_tables"`
}

type TestFramework struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Runner  string `json:"runner"`
}

type CIConfig struct {
	Platform string `json:"platform"`
	Config   string `json:"config"`
}

type Scores struct {
	RealityConformance  float64 `json:"reality_conformance"`
	DependencyIntegrity float64 `json:"dependency_integrity"`
	StructuralAlignment float64 `json:"structural_alignment"`
}

type SnapshotData struct {
	ID            uuid.UUID            `json:"id"`
	ProjectID     uuid.UUID            `json:"project_id"`
	RoadmapItemID uuid.UUID            `json:"roadmap_item_id"`
	State         SnapshotState        `json:"state"`
	SnapshotJSON  *EnvironmentSnapshot `json:"snapshot_json,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

type SnapshotAnalysis struct {
	ID                 uuid.UUID              `json:"id"`
	SnapshotID         uuid.UUID              `json:"snapshot_id"`
	Scores             Scores                 `json:"scores"`
	Verdict            string                 `json:"verdict"`
	DriftDetected      bool                   `json:"drift_detected"`
	AlignmentConflicts map[string]interface{} `json:"alignment_conflicts"`
	CreatedAt          time.Time              `json:"created_at"`
}
