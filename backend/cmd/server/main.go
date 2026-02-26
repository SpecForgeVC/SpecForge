package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/scott/specforge/internal/api"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/drift"
	"github.com/scott/specforge/internal/infra"
	"github.com/scott/specforge/internal/infra/auth"
	"github.com/scott/specforge/internal/infra/db"
	"github.com/scott/specforge/internal/logger"
	"github.com/scott/specforge/internal/mcp"
	mw "github.com/scott/specforge/internal/transport/middleware"
	"github.com/scott/specforge/internal/ui_roadmap"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	e := echo.New()

	// Middleware
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// Auth Initialization
	authCfg := auth.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		Issuer:     os.Getenv("JWT_ISSUER"),
		Audience:   os.Getenv("JWT_AUDIENCE"),
		Algorithm:  os.Getenv("JWT_ALGORITHM"), // Defaults to HS256 in validator if empty
	}
	validator := auth.NewJWTValidator(authCfg)

	// Adapters for Echo
	authMiddleware := mw.Adapt(mw.AuthMiddleware(validator))
	requireRole := func(roles ...domain.Role) echo.MiddlewareFunc {
		return mw.Adapt(mw.RequireRole(roles...))
	}

	// Database connection (Postgres)
	dbConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	// Repositories
	userRepo := infra.NewUserRepository(queries)
	wsRepo := infra.NewWorkspaceRepository(queries)
	pRepo := infra.NewProjectRepository(queries)
	rmRepo := infra.NewRoadmapItemRepository(queries)
	cRepo := infra.NewContractRepository(queries)
	sRepo := infra.NewSnapshotRepository(queries)
	propRepo := infra.NewAiProposalRepository(queries)
	auditRepo := infra.NewAuditLogRepository(queries)
	reqRepo := infra.NewRequirementRepository(queries)
	varRepo := infra.NewVariableRepository(queries)
	whRepo := infra.NewWebhookRepository(queries)
	valRepo := infra.NewValidationRuleRepository(queries)
	llmRepo := infra.NewLLMRepository(dbConn)
	fiRepo := infra.NewFeatureIntelligenceRepository(queries)
	vlRepo := infra.NewVariableLineageRepository(queries)
	refRepo := infra.NewRefinementRepository(dbConn)
	sessionRepo := infra.NewImportSessionRepository(dbConn)

	// Intelligence Alignment Repos (using sql.DB for now)
	alignmentRepo := infra.NewAlignmentRepository(dbConn)
	depRepo := infra.NewRoadmapDependencyRepository(dbConn)

	diffEngine := drift.NewDiffEngine()

	// Services
	authService := app.NewAuthService(userRepo, string(authCfg.SigningKey), authCfg.Issuer, authCfg.Audience)
	auditService := app.NewAuditLogService(auditRepo)

	// Drift
	driftService := drift.NewDriftService(cRepo, sRepo, diffEngine, auditService)

	// Notifications
	notifyService := app.NewNotificationService()

	// Initialize Intelligence Service before others that depend on it
	fiService := app.NewFeatureIntelligenceService(fiRepo, rmRepo, cRepo, varRepo, reqRepo, driftService, notifyService)
	vlService := app.NewVariableLineageService(vlRepo)

	// NEW: Alignment & Dependency Services
	alignmentService := app.NewAlignmentService(alignmentRepo, rmRepo, depRepo, cRepo, varRepo, valRepo)
	depService := app.NewRoadmapDependencyService(depRepo, rmRepo, auditService)

	govService := app.NewGovernanceService(fiService, propRepo, varRepo)

	wsService := app.NewWorkspaceService(wsRepo, auditService)
	pService := app.NewProjectService(pRepo, auditService)
	rmService := app.NewRoadmapItemService(depRepo, rmRepo, auditService, fiService, govService, alignmentService)
	cService := app.NewContractService(cRepo, rmRepo, fiService, govService, alignmentService)
	sService := app.NewSnapshotService(sRepo)
	propService := app.NewAiProposalService(propRepo, rmRepo, sRepo, diffEngine, auditService)
	reqService := app.NewRequirementService(reqRepo, auditService)
	varService := app.NewVariableService(varRepo, cRepo, rmRepo, auditService, fiService, alignmentService)
	whService := app.NewWebhookService(whRepo, auditService)
	valService := app.NewValidationRuleService(valRepo, auditService)

	// LLM & Refinement
	llmFactory := infra.NewLLMFactory()
	llmService := app.NewLLMService(llmRepo, llmFactory)
	refService := app.NewRefinementService(refRepo, pRepo, llmService)

	// Bootstrap Intelligence
	bootstrapRepo := infra.NewBootstrapRepository(queries)
	bootstrapService := app.NewBootstrapService(bootstrapRepo, pRepo, sessionRepo)

	// Build Artifact Export
	artifactExporter := infra.NewArtifactExporter()
	artifactService := app.NewBuildArtifactService(rmRepo, cRepo, varRepo, reqRepo, valRepo, govService)

	// UI Roadmap Engine
	uiRoadmapRepo := ui_roadmap.NewRepository(dbConn)
	uiRoadmapService := ui_roadmap.NewService(uiRoadmapRepo)
	uiRoadmapHandler := api.NewUIRoadmapHandler(uiRoadmapService)

	// MCP Token System
	mcpTokenRepo := infra.NewMCPTokenRepository(dbConn)
	mcpTokenService := app.NewMCPTokenService(mcpTokenRepo)

	// MCP Server Integration
	mcpRepo := infra.NewMCPRepository(dbConn)
	importService := app.NewImportService(pRepo, mcpRepo, alignmentService, app.NewDiffService(), bootstrapRepo, sessionRepo)
	mcpHandlers := mcp.NewHandlers(mcpRepo, importService)
	mcpConfig := mcp.Config{
		Port:         8081,
		BindAddress:  "0.0.0.0",
		Enabled:      true,
		AuthRequired: true,
		Token:        os.Getenv("MCP_TOKEN"),
		TokenService: mcpTokenService, // Pass token service for validation
	}
	if mcpConfig.Token == "" {
		mcpConfig.Token = "default-rae-token-change-me"
	}
	mcpServer := mcp.NewServer(mcpConfig, mcpHandlers)
	if err := mcpServer.Start(); err != nil {
		logger.Log.Error("Failed to start MCP Server", zap.Error(err))
	}

	// Handlers
	authHandler := api.NewAuthHandler(authService)
	wsHandler := api.NewWorkspaceHandler(wsService)
	pHandler := api.NewProjectHandler(pService)
	rmHandler := api.NewRoadmapItemHandler(rmService, artifactService, artifactExporter)
	cHandler := api.NewContractHandler(cService)
	sHandler := api.NewSnapshotHandler(sService)
	propHandler := api.NewAiProposalHandler(propService)
	auditHandler := api.NewAuditLogHandler(auditService)
	reqHandler := api.NewRequirementHandler(reqService)
	varHandler := api.NewVariableHandler(varService)
	whHandler := api.NewWebhookHandler(whService)
	valHandler := api.NewValidationRuleHandler(valService)
	driftHandler := api.NewDriftHandler(driftService)
	fiHandler := api.NewFeatureIntelligenceHandler(fiService)
	vlHandler := api.NewVariableLineageHandler(vlService)
	webSocketHandler := api.NewWSHandler(notifyService, validator)
	llmHandler := api.NewLLMSettingsHandler(llmService)
	refHandler := api.NewRefinementHandler(refService)
	bootstrapHandler := api.NewBootstrapHandler(bootstrapService)
	mcpHandler := api.NewMCPHandler(mcpTokenService, pService, mcpServer)

	// NEW: Alignment & Dependency Handlers
	alignmentHandler := api.NewAlignmentHandler(alignmentService)
	depHandler := api.NewRoadmapDependencyHandler(depService)

	// Routes
	v1 := e.Group("/api/v1")

	// Public Routes
	v1.POST("/auth/login", authHandler.Login)
	v1.POST("/auth/refresh", authHandler.Refresh)
	v1.GET("/ws", webSocketHandler.Connect)
	v1.GET("/auth/me", authHandler.GetMe, authMiddleware)

	// Protected Routes (everything below this point requires authentication)
	v1.Use(authMiddleware)

	// Audit Logs
	v1.GET("/audit-logs/:entityType/:entityId", auditHandler.GetEntityLogs, requireRole(domain.RoleReviewer, domain.RoleAdmin, domain.RoleOwner))

	// Workspaces
	v1.GET("/workspaces", wsHandler.ListWorkspaces)
	v1.POST("/workspaces", wsHandler.CreateWorkspace, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.GET("/workspaces/:workspaceId", wsHandler.GetWorkspace)
	v1.PATCH("/workspaces/:workspaceId", wsHandler.UpdateWorkspace, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.DELETE("/workspaces/:workspaceId", wsHandler.DeleteWorkspace, requireRole(domain.RoleOwner))

	// Projects
	v1.GET("/workspaces/:workspaceId/projects", pHandler.ListProjects)
	v1.POST("/workspaces/:workspaceId/projects", pHandler.CreateProject, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/projects/:projectId", pHandler.GetProject)
	v1.PATCH("/projects/:projectId", pHandler.UpdateProject, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/projects/:projectId", pHandler.DeleteProject, requireRole(domain.RoleOwner, domain.RoleAdmin))

	// Roadmap Items
	v1.GET("/projects/:projectId/roadmap-items", rmHandler.ListRoadmapItems)
	v1.POST("/projects/:projectId/roadmap-items", rmHandler.CreateRoadmapItem, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleReviewer))
	v1.GET("/projects/:projectId/contracts", cHandler.ListContractsByProject)
	v1.POST("/projects/:projectId/contracts", cHandler.CreateContractByProject, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/projects/:projectId/variables", varHandler.ListVariablesByProject)
	v1.POST("/projects/:projectId/variables", varHandler.CreateVariableByProject, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/projects/:projectId/snapshots", sHandler.ListSnapshotsByProject)
	v1.GET("/roadmap-items/:roadmapItemId", rmHandler.GetRoadmapItem)
	v1.PATCH("/roadmap-items/:roadmapItemId", rmHandler.UpdateRoadmapItem, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleReviewer))
	v1.DELETE("/roadmap-items/:roadmapItemId", rmHandler.DeleteRoadmapItem, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/roadmap-items/:roadmapItemId/export", rmHandler.ExportRoadmapItem)

	// Intelligence Alignment
	v1.GET("/projects/:projectId/alignment", alignmentHandler.GetAlignmentReport)
	v1.POST("/projects/:projectId/alignment", alignmentHandler.TriggerAlignmentCheck, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/projects/:projectId/roadmap-dependencies", depHandler.ListDependencies)
	v1.POST("/projects/:projectId/roadmap-dependencies", depHandler.CreateDependency, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/roadmap-dependencies/:dependencyId", depHandler.DeleteDependency, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))

	// AI Proposals
	v1.GET("/projects/:projectId/ai-proposals", propHandler.ListProposals)
	v1.POST("/ai-proposals", propHandler.CreateProposal, requireRole(domain.RoleAIAgent, domain.RoleEngineer, domain.RoleOwner, domain.RoleAdmin))
	v1.GET("/ai-proposals/:proposalId", propHandler.GetProposal)
	v1.POST("/ai-proposals/:proposalId/approve", propHandler.ApproveProposal, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleReviewer))
	v1.POST("/ai-proposals/:proposalId/reject", propHandler.RejectProposal, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleReviewer))

	// Contracts
	v1.GET("/roadmap-items/:roadmapItemId/contracts", cHandler.ListContracts)
	v1.POST("/roadmap-items/:roadmapItemId/contracts", cHandler.CreateContract)
	v1.GET("/contracts/:contractId", cHandler.GetContract)
	v1.PATCH("/contracts/:contractId", cHandler.UpdateContract, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/contracts/:contractId", cHandler.DeleteContract)

	// Snapshots
	v1.GET("/roadmap-items/:roadmapItemId/snapshots", sHandler.ListSnapshots)
	v1.POST("/roadmap-items/:roadmapItemId/snapshots", sHandler.CreateSnapshot, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleAIAgent))
	v1.GET("/snapshots/:snapshotId", sHandler.GetSnapshot)

	// Requirements
	v1.GET("/roadmap-items/:roadmapItemId/requirements", reqHandler.ListRequirements)
	v1.POST("/roadmap-items/:roadmapItemId/requirements", reqHandler.CreateRequirement, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/requirements/:requirementId", reqHandler.GetRequirement)
	v1.PATCH("/requirements/:requirementId", reqHandler.UpdateRequirement, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/requirements/:requirementId", reqHandler.DeleteRequirement, requireRole(domain.RoleOwner, domain.RoleAdmin))

	// Variables
	v1.GET("/contracts/:contractId/variables", varHandler.ListVariables)
	v1.POST("/contracts/:contractId/variables", varHandler.CreateVariable, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/variables/:variableId", varHandler.GetVariable)
	v1.PATCH("/variables/:variableId", varHandler.UpdateVariable, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/variables/:variableId", varHandler.DeleteVariable, requireRole(domain.RoleOwner, domain.RoleAdmin))

	// Intelligence & Lineage
	v1.GET("/roadmap-items/:roadmapItemId/intelligence", fiHandler.GetFeatureIntelligence)
	v1.GET("/variables/:variableId/events", vlHandler.GetLineageEvents)
	v1.GET("/variables/:variableId/lineage", vlHandler.GetLineageGraph)

	// Webhooks
	v1.GET("/projects/:projectId/webhooks", whHandler.ListWebhooks)
	v1.POST("/projects/:projectId/webhooks", whHandler.CreateWebhook, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.GET("/webhooks/:webhookId", whHandler.GetWebhook)
	v1.PATCH("/webhooks/:webhookId", whHandler.UpdateWebhook, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.DELETE("/webhooks/:webhookId", whHandler.DeleteWebhook, requireRole(domain.RoleOwner))

	// Validation Rules
	v1.GET("/projects/:projectId/validation-rules", valHandler.ListValidationRules)
	v1.POST("/projects/:projectId/validation-rules", valHandler.CreateValidationRule, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.GET("/validation-rules/:ruleId", valHandler.GetValidationRule)
	v1.PATCH("/validation-rules/:ruleId", valHandler.UpdateValidationRule, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.DELETE("/validation-rules/:ruleId", valHandler.DeleteValidationRule, requireRole(domain.RoleOwner))

	// Drift Check
	v1.POST("/contracts/:contractId/drift-check", driftHandler.RunDriftCheck, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleReviewer))
	v1.GET("/drift/history", driftHandler.GetDriftHistory, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleReviewer))
	v1.POST("/drift/generate-fixes", driftHandler.GenerateDriftFixes, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleReviewer))

	// Activity Logs
	v1.GET("/roadmap-items/:roadmapItemId/activity", auditHandler.GetRoadmapItemActivity)

	// LLM Settings
	v1.GET("/settings/llm", llmHandler.GetConfig, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.PUT("/settings/llm", llmHandler.UpdateConfig, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.POST("/settings/llm/test", llmHandler.TestConnection, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.GET("/settings/llm/warmup", llmHandler.Warmup, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.POST("/settings/llm/models", llmHandler.ListModels, requireRole(domain.RoleOwner, domain.RoleAdmin))

	// Refinement
	v1.POST("/refinement", refHandler.StartSession, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleAIAgent))
	v1.GET("/refinement/:sessionId/events", refHandler.StreamEvents, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer, domain.RoleAIAgent))
	v1.POST("/refinement/:sessionId/approve", refHandler.ApproveSession, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))

	// Bootstrap Intelligence
	v1.POST("/projects/:projectId/bootstrap/generate-prompt", bootstrapHandler.GeneratePrompt, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.POST("/projects/:projectId/bootstrap/ingest", bootstrapHandler.IngestBootstrap, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/projects/:projectId/bootstrap/snapshots", bootstrapHandler.ListSnapshots)
	v1.GET("/projects/:projectId/bootstrap/snapshots/:snapshotId", bootstrapHandler.GetSnapshot)
	v1.GET("/projects/:projectId/bootstrap/latest", bootstrapHandler.GetLatestSnapshot)
	v1.GET("/projects/:projectId/bootstrap/session", bootstrapHandler.GetLatestImportSession)
	v1.POST("/projects/:projectId/bootstrap/diff", bootstrapHandler.DiffSnapshots, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))

	// UI Roadmap Items
	v1.GET("/projects/:projectId/ui-roadmap", uiRoadmapHandler.ListItems)
	v1.POST("/projects/:projectId/ui-roadmap", uiRoadmapHandler.SaveItem, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/ui-roadmap/:id", uiRoadmapHandler.GetItem)
	v1.PUT("/ui-roadmap/:id", uiRoadmapHandler.SaveItem, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/ui-roadmap/:id", uiRoadmapHandler.DeleteItem, requireRole(domain.RoleOwner, domain.RoleAdmin))
	v1.GET("/ui-roadmap/:id/export", uiRoadmapHandler.ExportItem)
	v1.POST("/ui-roadmap/:id/sync", uiRoadmapHandler.SyncFigma)
	v1.GET("/ui-roadmap/:id/plugin-assets", uiRoadmapHandler.GetPluginAssets)

	// MCP Management
	v1.GET("/mcp/status", mcpHandler.GetStatus)
	v1.POST("/mcp/tokens", mcpHandler.GenerateToken, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/mcp/tokens", mcpHandler.ListTokens, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.DELETE("/mcp/tokens/:id", mcpHandler.RevokeToken, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))
	v1.GET("/mcp/config/download", mcpHandler.DownloadConfig, requireRole(domain.RoleOwner, domain.RoleAdmin, domain.RoleEngineer))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Shutting down the server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	if err := mcpServer.Stop(ctx); err != nil {
		logger.Log.Error("Failed to stop MCP Server", zap.Error(err))
	}
}
