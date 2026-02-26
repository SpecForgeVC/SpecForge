package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type RoadmapDependencyHandler struct {
	service app.RoadmapDependencyService
}

func NewRoadmapDependencyHandler(service app.RoadmapDependencyService) *RoadmapDependencyHandler {
	return &RoadmapDependencyHandler{service: service}
}

type roadmapDependencyCreateRequest struct {
	SourceID       uuid.UUID             `json:"source_id"`
	TargetID       uuid.UUID             `json:"target_id"`
	DependencyType domain.DependencyType `json:"dependency_type"`
}

func (h *RoadmapDependencyHandler) CreateDependency(c echo.Context) error {
	var req roadmapDependencyCreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	dep, err := h.service.CreateDependency(c.Request().Context(), req.SourceID, req.TargetID, req.DependencyType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, dep)
}

func (h *RoadmapDependencyHandler) ListDependencies(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}

	deps, err := h.service.ListDependencies(c.Request().Context(), projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    deps,
	})
}

func (h *RoadmapDependencyHandler) DeleteDependency(c echo.Context) error {
	depID, err := uuid.Parse(c.Param("dependencyId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dependency ID")
	}

	if err := h.service.DeleteDependency(c.Request().Context(), depID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
