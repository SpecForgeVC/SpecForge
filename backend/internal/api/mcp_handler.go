package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/logger"
	"github.com/scott/specforge/internal/mcp"
	"github.com/scott/specforge/internal/transport/middleware"
	"go.uber.org/zap"
)

type MCPHandler struct {
	tokenService   app.MCPTokenService
	projectService app.ProjectService
	mcpServer      *mcp.Server
}

func NewMCPHandler(tokenService app.MCPTokenService, projectService app.ProjectService, mcpServer *mcp.Server) *MCPHandler {
	return &MCPHandler{
		tokenService:   tokenService,
		projectService: projectService,
		mcpServer:      mcpServer,
	}
}

type MCPStatusResponse struct {
	Enabled      bool   `json:"enabled"`
	Port         int    `json:"port"`
	BindAddress  string `json:"bind_address"`
	IsRunning    bool   `json:"is_running"`
	AuthRequired bool   `json:"auth_required"`
}

func (h *MCPHandler) GetStatus(c echo.Context) error {
	cfg := h.mcpServer.GetConfig()
	return c.JSON(http.StatusOK, MCPStatusResponse{
		Enabled:      cfg.Enabled,
		Port:         cfg.Port,
		BindAddress:  cfg.BindAddress,
		IsRunning:    h.mcpServer.IsRunning(),
		AuthRequired: cfg.AuthRequired,
	})
}

func (h *MCPHandler) GenerateToken(c echo.Context) error {
	projectIDStr := strings.TrimSpace(c.QueryParam("project_id"))
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid project_id")
	}

	principal, ok := middleware.PrincipalFromContext(c.Request().Context())
	if !ok {
		logger.Error("GenerateToken: failed to get principal from context")
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	logger.Info("GenerateToken Request",
		zap.String("user_id", principal.UserID.String()),
		zap.String("project_id", projectID.String()),
	)

	// Verify project ownership (simplified)
	p, err := h.projectService.GetProject(c.Request().Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "project not found")
		}
		logger.Log.Error("GenerateToken: GetProject failed",
			zap.String("project_id", projectID.String()),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	logger.Info("GenerateToken: project verified", zap.String("project_id", p.ID.String()))

	resp, err := h.tokenService.GenerateToken(c.Request().Context(), principal.UserID, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *MCPHandler) ListTokens(c echo.Context) error {
	projectIDStr := strings.TrimSpace(c.QueryParam("project_id"))
	logger.Info("ListTokens Request", zap.String("project_id", projectIDStr))
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid project_id")
	}

	tokens, err := h.tokenService.ListTokens(c.Request().Context(), projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list tokens")
	}

	return c.JSON(http.StatusOK, tokens)
}

func (h *MCPHandler) RevokeToken(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid token id")
	}

	if err := h.tokenService.RevokeToken(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to revoke token")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *MCPHandler) DownloadConfig(c echo.Context) error {
	ide := c.QueryParam("ide")
	projectIDStr := strings.TrimSpace(c.QueryParam("project_id"))
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid project_id")
	}

	if _, ok := middleware.PrincipalFromContext(c.Request().Context()); !ok {
		logger.Error("DownloadConfig: failed to get principal from context")
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	logger.Info("DownloadConfig Request",
		zap.String("ide", ide),
		zap.String("project_id", projectID.String()),
	)

	// Verify project ownership
	if _, err := h.projectService.GetProject(c.Request().Context(), projectID); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "project not found")
		}
		logger.Log.Error("DownloadConfig: GetProject failed",
			zap.String("project_id", projectID.String()),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	// Get latest token for this project (simplified)
	tokens, err := h.tokenService.ListTokens(c.Request().Context(), projectID)
	if err != nil || len(tokens) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no active token found for project")
	}

	// Find first non-revoked token
	var activeToken string
	for _, t := range tokens {
		if !t.Revoked {
			// We can't get the raw token back if it's hashed.
			// HMM. The requirement says "dynamic config generator".
			// If we only store the hash, we can't show it in the config file download UNLESS we generate it right then.
			// OR we ask the user to provide the token they copied earlier.
			// Re-reading requirements: "Step 5 - Download Config".
			// "GET /api/v1/mcp/config/download?ide=cursor&project_id=uuid".
			// If we want the token in the file, we either need to keep it in plain text (bad) OR generate a NEW one during download (maybe?)
			// OR the UI provides it in the download request?
			// Let's assume the UI passes the token it JUST generated, or we generated it on the fly.
			// BUT if they refresh the page, they can't download it again with the same token.
			// This is a trade-off. I'll make it so if they don't provide a token in the request, we use a placeholder or they must provide it.
			// Actually, let's allow the UI to pass the token in the query param for the download link.
			activeToken = c.QueryParam("token")
			break
		}
	}

	if activeToken == "" {
		activeToken = "REPLACE_WITH_YOUR_TOKEN"
	}

	mcpUrl := fmt.Sprintf("http://localhost:%d", h.mcpServer.GetConfig().Port)
	filename := fmt.Sprintf("specforge-%s-%s.json", ide, projectIDStr[:8])

	switch ide {
	case "cursor":
		config := map[string]interface{}{
			"specforge": map[string]interface{}{
				"mcpServerUrl": mcpUrl,
				"projectId":    projectIDStr,
				"apiToken":     activeToken,
				"autoConnect":  true,
				"cliPath":      "specforge-mcp",
			},
		}
		c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
		return c.JSON(http.StatusOK, config)
	case "anti-gravity":
		config := map[string]interface{}{
			"mcpServers": map[string]interface{}{
				"specforge": map[string]interface{}{
					"serverUrl": mcpUrl,
					"headers": map[string]interface{}{
						"Authorization": "Bearer " + activeToken,
						"Content-Type":  "application/json",
					},
				},
			},
		}
		c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
		return c.JSON(http.StatusOK, config)
	case "claude":
		config := map[string]interface{}{
			"tools": map[string]interface{}{
				"specforge": map[string]interface{}{
					"command": "specforge-mcp",
					"args": []string{
						"--server", mcpUrl,
						"--token", activeToken,
						"--project", projectIDStr,
					},
				},
			},
		}
		c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", filename))
		return c.JSON(http.StatusOK, config)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "unsupported ide")
	}
}
