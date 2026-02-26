package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	mw "github.com/scott/specforge/internal/transport/middleware"
)

type ProjectHandler struct {
	service app.ProjectService
}

func NewProjectHandler(service app.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) GetProject(c echo.Context) error {
	id, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	p, err := h.service.GetProject(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "project not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, p)
}

func (h *ProjectHandler) ListProjects(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid workspace id", "")
	}
	projects, err := h.service.ListProjects(c.Request().Context(), workspaceID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list projects", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, projects)
}

type projectCreateRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	TechStack     map[string]interface{} `json:"tech_stack"`
	Settings      map[string]interface{} `json:"settings"`
	MCPSettings   domain.MCPSettings     `json:"mcp_settings"`
	RepositoryURL string                 `json:"repository_url"`
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid workspace id", "")
	}
	req := new(projectCreateRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "failed to bind request", err.Error())
	}

	userID := GetUserID(c)

	p, err := h.service.CreateProject(c.Request().Context(), workspaceID, req.Name, req.Description, req.TechStack, req.Settings, req.MCPSettings, req.RepositoryURL, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create project", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, p)
}

func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	id, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	req := new(projectCreateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	p, err := h.service.UpdateProject(c.Request().Context(), id, req.Name, req.Description, req.TechStack, req.Settings, req.MCPSettings, req.RepositoryURL, principal.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	id, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if err := h.service.DeleteProject(c.Request().Context(), id, principal.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
