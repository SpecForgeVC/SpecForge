package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/scott/specforge/internal/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// --- Repositories ---

type WorkspaceRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Workspace, error)
	List(ctx context.Context) ([]domain.Workspace, error)
	Create(ctx context.Context, ws *domain.Workspace) error
	Update(ctx context.Context, ws *domain.Workspace) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	List(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error)
	Create(ctx context.Context, p *domain.Project) error
	Update(ctx context.Context, p *domain.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RoadmapItemRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.RoadmapItem, error)
	List(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapItem, error)
	Create(ctx context.Context, item *domain.RoadmapItem) error
	Update(ctx context.Context, item *domain.RoadmapItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ContractRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error)
	List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.ContractDefinition, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.ContractDefinition, error)
	Create(ctx context.Context, c *domain.ContractDefinition) error
	Update(ctx context.Context, c *domain.ContractDefinition) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SnapshotRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.VersionSnapshot, error)
	List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.VersionSnapshot, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VersionSnapshot, error)
	Create(ctx context.Context, s *domain.VersionSnapshot) error
}

// --- Services ---

type WorkspaceService interface {
	GetWorkspace(ctx context.Context, id uuid.UUID) (*domain.Workspace, error)
	ListWorkspaces(ctx context.Context) ([]domain.Workspace, error)
	CreateWorkspace(ctx context.Context, name, description string, userID uuid.UUID) (*domain.Workspace, error)
	UpdateWorkspace(ctx context.Context, id uuid.UUID, name, description string, userID uuid.UUID) (*domain.Workspace, error)
	DeleteWorkspace(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type ProjectService interface {
	GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	ListProjects(ctx context.Context, workspaceID uuid.UUID) ([]domain.Project, error)
	CreateProject(ctx context.Context, workspaceID uuid.UUID, name, description string, techStack map[string]interface{}, settings map[string]interface{}, mcpSettings domain.MCPSettings, repositoryURL string, userID uuid.UUID) (*domain.Project, error)
	UpdateProject(ctx context.Context, id uuid.UUID, name, description string, techStack map[string]interface{}, settings map[string]interface{}, mcpSettings domain.MCPSettings, repositoryURL string, userID uuid.UUID) (*domain.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	RecommendStack(ctx context.Context, purpose, constraints string) (*domain.TechStackRecommendation, error)
}

type RoadmapItemService interface {
	GetRoadmapItem(ctx context.Context, id uuid.UUID) (*domain.RoadmapItem, error)
	ListRoadmapItems(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapItem, error)
	CreateRoadmapItem(ctx context.Context, item *domain.RoadmapItem, userID uuid.UUID) (*domain.RoadmapItem, error)
	UpdateRoadmapItem(ctx context.Context, id uuid.UUID, title, description, businessContext, technicalContext string, status domain.RoadmapItemStatus, userID uuid.UUID) (*domain.RoadmapItem, error)
	DeleteRoadmapItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type ContractService interface {
	GetContract(ctx context.Context, id uuid.UUID) (*domain.ContractDefinition, error)
	ListContracts(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.ContractDefinition, error)
	ListContractsByProject(ctx context.Context, projectID uuid.UUID) ([]domain.ContractDefinition, error)
	CreateContract(ctx context.Context, roadmapItemID uuid.UUID, cType domain.ContractType, version string, input, output, errSchema map[string]interface{}) (*domain.ContractDefinition, error)
	UpdateContract(ctx context.Context, id uuid.UUID, cType domain.ContractType, version string, input, output, errSchema map[string]interface{}) (*domain.ContractDefinition, error)
	DeleteContract(ctx context.Context, id uuid.UUID) error
}

type SnapshotService interface {
	GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.VersionSnapshot, error)
	ListSnapshots(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.VersionSnapshot, error)
	ListSnapshotsByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VersionSnapshot, error)
	CreateSnapshot(ctx context.Context, roadmapItemID uuid.UUID, data map[string]interface{}, createdBy uuid.UUID) (*domain.VersionSnapshot, error)
}

type AiProposalRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.AiProposal, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.AiProposal, error)
	ListByRoadmapItem(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.AiProposal, error)
	Create(ctx context.Context, p *domain.AiProposal) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProposalStatus, reviewedBy uuid.UUID) error
}

type AiProposalService interface {
	GetProposal(ctx context.Context, id uuid.UUID) (*domain.AiProposal, error)
	ListProposals(ctx context.Context, projectID uuid.UUID) ([]domain.AiProposal, error)
	CreateProposal(ctx context.Context, roadmapItemID uuid.UUID, pType domain.ProposalType, diff map[string]interface{}, reasoning string, confidence float64) (*domain.AiProposal, error)
	ApproveProposal(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	RejectProposal(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.AuditLog, error)
	ListDriftEvents(ctx context.Context) ([]domain.AuditLog, error)
}

type AuditLogService interface {
	Log(ctx context.Context, entityType string, entityID uuid.UUID, action string, userID uuid.UUID, oldData, newData map[string]interface{}) error
	GetEntityLogs(ctx context.Context, entityType string, entityID uuid.UUID) ([]domain.AuditLog, error)
	ListDriftEvents(ctx context.Context) ([]domain.AuditLog, error)
}

// Requirements
type RequirementRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Requirement, error)
	List(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.Requirement, error)
	Create(ctx context.Context, r *domain.Requirement) error
	Update(ctx context.Context, r *domain.Requirement) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RequirementService interface {
	GetRequirement(ctx context.Context, id uuid.UUID) (*domain.Requirement, error)
	ListRequirements(ctx context.Context, roadmapItemID uuid.UUID) ([]domain.Requirement, error)
	CreateRequirement(ctx context.Context, roadmapItemID uuid.UUID, title, description string, testable bool, acceptanceCriteria string, orderIndex int, userID uuid.UUID) (*domain.Requirement, error)
	UpdateRequirement(ctx context.Context, id uuid.UUID, title, description string, testable bool, acceptanceCriteria string, orderIndex int, userID uuid.UUID) (*domain.Requirement, error)
	DeleteRequirement(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// Variables
type VariableRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.VariableDefinition, error)
	List(ctx context.Context, contractID uuid.UUID) ([]domain.VariableDefinition, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VariableDefinition, error)
	Create(ctx context.Context, v *domain.VariableDefinition) error
	Update(ctx context.Context, v *domain.VariableDefinition) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VariableService interface {
	GetVariable(ctx context.Context, id uuid.UUID) (*domain.VariableDefinition, error)
	ListVariables(ctx context.Context, contractID uuid.UUID) ([]domain.VariableDefinition, error)
	ListVariablesByProject(ctx context.Context, projectID uuid.UUID) ([]domain.VariableDefinition, error)
	CreateVariable(ctx context.Context, contractID uuid.UUID, name, vType string, required bool, defaultValue, description string, validationRules map[string]interface{}, userID uuid.UUID) (*domain.VariableDefinition, error)
	UpdateVariable(ctx context.Context, id uuid.UUID, name, vType string, required bool, defaultValue, description string, validationRules map[string]interface{}, userID uuid.UUID) (*domain.VariableDefinition, error)
	DeleteVariable(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// Webhooks
type WebhookRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Webhook, error)
	List(ctx context.Context, projectID uuid.UUID) ([]domain.Webhook, error)
	Create(ctx context.Context, w *domain.Webhook) error
	Update(ctx context.Context, w *domain.Webhook) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WebhookService interface {
	GetWebhook(ctx context.Context, id uuid.UUID) (*domain.Webhook, error)
	ListWebhooks(ctx context.Context, projectID uuid.UUID) ([]domain.Webhook, error)
	CreateWebhook(ctx context.Context, projectID uuid.UUID, url string, events []string, secret string, userID uuid.UUID) (*domain.Webhook, error)
	UpdateWebhook(ctx context.Context, id uuid.UUID, url string, events []string, secret string, active bool, userID uuid.UUID) (*domain.Webhook, error)
	DeleteWebhook(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// Validation Rules
type ValidationRuleRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.ValidationRule, error)
	List(ctx context.Context, projectID uuid.UUID) ([]domain.ValidationRule, error)
	Create(ctx context.Context, r *domain.ValidationRule) error
	Update(ctx context.Context, r *domain.ValidationRule) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ValidationRuleService interface {
	GetValidationRule(ctx context.Context, id uuid.UUID) (*domain.ValidationRule, error)
	ListValidationRules(ctx context.Context, projectID uuid.UUID) ([]domain.ValidationRule, error)
	CreateValidationRule(ctx context.Context, projectID uuid.UUID, name, rType string, config map[string]interface{}, description string, userID uuid.UUID) (*domain.ValidationRule, error)
	UpdateValidationRule(ctx context.Context, id uuid.UUID, name, rType string, config map[string]interface{}, description string, userID uuid.UUID) (*domain.ValidationRule, error)
	DeleteValidationRule(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// Feature Intelligence
type FeatureIntelligenceRepository interface {
	Get(ctx context.Context, featureID uuid.UUID) (*domain.FeatureIntelligence, error)
	Create(ctx context.Context, fi *domain.FeatureIntelligence) error
	Update(ctx context.Context, fi *domain.FeatureIntelligence) error
	Delete(ctx context.Context, featureID uuid.UUID) error
}

type FeatureIntelligenceService interface {
	GetFeatureScore(ctx context.Context, featureID uuid.UUID) (*domain.FeatureIntelligence, error)
	CalculateFeatureScore(ctx context.Context, featureID uuid.UUID) (*domain.FeatureIntelligence, error)
}

// Variable Lineage
type VariableLineageRepository interface {
	CreateEvent(ctx context.Context, event *domain.VariableLineageEvent) error
	ListEvents(ctx context.Context, variableID uuid.UUID) ([]domain.VariableLineageEvent, error)
	CreateDependency(ctx context.Context, dep *domain.VariableDependency) error
	ListDependencies(ctx context.Context, variableID uuid.UUID) ([]domain.VariableDependency, error)
}

type VariableLineageService interface {
	TrackEvent(ctx context.Context, variableID uuid.UUID, eventType domain.LineageEventType, source, description string, userID uuid.UUID, metadata map[string]interface{}) error
	GetLineageEvents(ctx context.Context, variableID uuid.UUID) ([]domain.VariableLineageEvent, error)
	GetLineageGraph(ctx context.Context, variableID uuid.UUID) (map[string]interface{}, error)
}

type GovernanceService interface {
	CanBuildFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error)
	CanDeployFeature(ctx context.Context, featureID uuid.UUID) (bool, []string, error)
	CanUpdateContract(ctx context.Context, contractID uuid.UUID) (bool, []string, error)
}

type DiffService interface {
	CompareSnapshots(ctx context.Context, oldSnap, newSnap map[string]interface{}) (*domain.DriftReport, error)
	CompareProjectSnapshot(ctx context.Context, snapshot domain.ProjectSnapshot, projectID uuid.UUID) (*domain.DriftReport, error)
}

// LLM
type LLMRepository interface {
	GetActive(ctx context.Context) (*domain.LLMConfiguration, error)
	Upsert(ctx context.Context, config *domain.LLMConfiguration) error
}

type LLMService interface {
	GetActiveConfig(ctx context.Context) (*domain.LLMConfiguration, error)
	UpdateConfig(ctx context.Context, config *domain.LLMConfiguration) error
	GetClient(ctx context.Context) (domain.LLMClient, error)
	TestConfiguration(ctx context.Context, config *domain.LLMConfiguration) error
	ListModels(ctx context.Context, config *domain.LLMConfiguration) ([]string, error)
}

// Refinement
type RefinementRepository interface {
	CreateSession(ctx context.Context, session *domain.RefinementSession) error
	GetSession(ctx context.Context, id uuid.UUID) (*domain.RefinementSession, error)
	UpdateSession(ctx context.Context, session *domain.RefinementSession) error
	CreateIteration(ctx context.Context, iteration *domain.RefinementIteration) error
}

type RefinementService interface {
	StartSession(ctx context.Context, artifactType string, targetType string, prompt string, contextData map[string]any, maxIterations int) (*domain.RefinementSession, error)
	GetSession(ctx context.Context, id uuid.UUID) (*domain.RefinementSession, error)
	GetSessionEvents(ctx context.Context, id uuid.UUID) (<-chan domain.RefinementEvent, error)
	ApproveSession(ctx context.Context, id uuid.UUID) error
}

// Project Bootstrap Intelligence
type BootstrapRepository interface {
	InsertSnapshot(ctx context.Context, s *domain.ProjectIntelligenceSnapshot) error
	ListSnapshots(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectIntelligenceSnapshot, error)
	GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error)
	GetLatestSnapshot(ctx context.Context, projectID uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error)
	GetMaxVersion(ctx context.Context, projectID uuid.UUID) (int, error)
	InsertModule(ctx context.Context, m *domain.ProjectModule) error
	InsertEntity(ctx context.Context, e *domain.ProjectEntity) error
	InsertApiEntry(ctx context.Context, a *domain.ProjectApiEntry) error
	InsertContractEntry(ctx context.Context, c *domain.ProjectContractEntry) error
	ListModulesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectModule, error)
	ListEntitiesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectEntity, error)
	ListApiEntriesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectApiEntry, error)
	ListContractEntriesBySnapshot(ctx context.Context, snapshotID uuid.UUID) ([]domain.ProjectContractEntry, error)
}

type BootstrapService interface {
	GeneratePrompt(ctx context.Context, projectID uuid.UUID) (string, string, error)
	IngestBootstrap(ctx context.Context, projectID uuid.UUID, payload domain.BootstrapPayload) (*domain.ProjectIntelligenceSnapshot, []string, error)
	ListSnapshots(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectIntelligenceSnapshot, error)
	GetSnapshot(ctx context.Context, snapshotID uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error)
	GetLatestSnapshot(ctx context.Context, projectID uuid.UUID) (*domain.ProjectIntelligenceSnapshot, error)
	GetLatestImportSession(ctx context.Context, projectID uuid.UUID) (*domain.ImportSession, error)
	DiffSnapshots(ctx context.Context, projectID uuid.UUID, fromID, toID *uuid.UUID) (map[string]interface{}, error)
}

// --- Intelligence Alignment ---

type RoadmapDependencyRepository interface {
	Create(ctx context.Context, dep *domain.RoadmapDependency) error
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapDependency, error)
	ListBySource(ctx context.Context, sourceID uuid.UUID) ([]domain.RoadmapDependency, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type AlignmentRepository interface {
	GetLatestReport(ctx context.Context, projectID uuid.UUID) (*domain.AlignmentReport, error)
	CreateReport(ctx context.Context, report *domain.AlignmentReport) error
}

type RoadmapDependencyService interface {
	CreateDependency(ctx context.Context, sourceID, targetID uuid.UUID, dType domain.DependencyType) (*domain.RoadmapDependency, error)
	ListDependencies(ctx context.Context, projectID uuid.UUID) ([]domain.RoadmapDependency, error)
	DeleteDependency(ctx context.Context, id uuid.UUID) error
}

type AlignmentService interface {
	GetAlignmentReport(ctx context.Context, projectID uuid.UUID) (*domain.AlignmentReport, error)
	TriggerAlignmentCheck(ctx context.Context, projectID uuid.UUID) (*domain.AlignmentReport, error)
	AnalyzeSnapshot(ctx context.Context, snapshot domain.ProjectSnapshot) (*domain.AlignmentReport, error)
}

type MCPRepository interface {
	CreateSnapshot(ctx context.Context, projectID, roadmapItemID uuid.UUID) (uuid.UUID, error)
	UpdateSnapshotState(ctx context.Context, id uuid.UUID, state domain.SnapshotState) error
	SaveSnapshotData(ctx context.Context, id uuid.UUID, data domain.EnvironmentSnapshot) error
	GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.SnapshotData, error)
	SaveAnalysis(ctx context.Context, analysis domain.SnapshotAnalysis) error
	ListActiveSnapshots(ctx context.Context, projectID uuid.UUID) ([]domain.SnapshotData, error)
}

type MCPTokenRepository interface {
	Create(ctx context.Context, token *domain.MCPToken) error
	GetByHash(ctx context.Context, hash string) (*domain.MCPToken, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]domain.MCPToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeAllForProject(ctx context.Context, projectID uuid.UUID) error
	UpdateUsage(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ImportService interface {
	InitializeImport(ctx context.Context, input domain.ToolInputInitProjectImport) (*domain.ToolOutputInitProjectImport, error)
	ProcessSnapshot(ctx context.Context, input domain.ToolInputSubmitProjectSnapshot) (*domain.ToolOutputSubmitProjectSnapshot, error)
	GetAlignmentRules(ctx context.Context, projectID uuid.UUID) (*domain.ToolOutputGetImportAlignmentRules, error)
	ProcessPostSnapshot(ctx context.Context, input domain.ToolInputSubmitPostImportSnapshot) (*domain.ToolOutputSubmitPostImportSnapshot, error)
	FinalizeImport(ctx context.Context, input domain.ToolInputFinalizeProjectImport) (*domain.ToolOutputFinalizeProjectImport, error)
}

type ImportSessionRepository interface {
	CreateSession(ctx context.Context, session *domain.ImportSession) error
	GetSession(ctx context.Context, id uuid.UUID) (*domain.ImportSession, error)
	GetLatestSessionByProject(ctx context.Context, projectID uuid.UUID) (*domain.ImportSession, error)
	UpdateSession(ctx context.Context, session *domain.ImportSession) error
	CreateArtifact(ctx context.Context, artifact *domain.ImportArtifact) error
	ListArtifactsBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.ImportArtifact, error)
}
