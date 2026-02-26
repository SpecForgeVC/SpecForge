package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type WorkspaceHandler struct {
	service app.WorkspaceService
}

func NewWorkspaceHandler(service app.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{service: service}
}

func (h *WorkspaceHandler) GetWorkspace(c echo.Context) error {
	id, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid workspace id", err.Error())
	}
	ws, err := h.service.GetWorkspace(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "workspace not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, ws)
}

func (h *WorkspaceHandler) ListWorkspaces(c echo.Context) error {
	pag := GetPagination(c)
	workspaces, err := h.service.ListWorkspaces(c.Request().Context())
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list workspaces", err.Error())
	}
	return SuccessResponseWithMeta(c, http.StatusOK, workspaces, pag)
}

type workspaceCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	req := new(workspaceCreateRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "failed to bind request", err.Error())
	}

	userID := GetUserID(c)

	ws, err := h.service.CreateWorkspace(c.Request().Context(), req.Name, req.Description, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create workspace", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) UpdateWorkspace(c echo.Context) error {
	id, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid workspace id", err.Error())
	}
	req := new(workspaceCreateRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "failed to bind request", err.Error())
	}

	userID := GetUserID(c)

	ws, err := h.service.UpdateWorkspace(c.Request().Context(), id, req.Name, req.Description, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update workspace", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, ws)
}

func (h *WorkspaceHandler) DeleteWorkspace(c echo.Context) error {
	id, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid workspace id", err.Error())
	}

	userID := GetUserID(c)

	if err := h.service.DeleteWorkspace(c.Request().Context(), id, userID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete workspace", err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
