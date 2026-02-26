package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/infra"
	mw "github.com/scott/specforge/internal/transport/middleware"
)

type RoadmapItemHandler struct {
	service         app.RoadmapItemService
	artifactService app.ArtifactService
	exporter        infra.ArtifactExporter
}

func NewRoadmapItemHandler(service app.RoadmapItemService, artifactService app.ArtifactService, exporter infra.ArtifactExporter) *RoadmapItemHandler {
	return &RoadmapItemHandler{
		service:         service,
		artifactService: artifactService,
		exporter:        exporter,
	}
}

func (h *RoadmapItemHandler) GetRoadmapItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}
	item, err := h.service.GetRoadmapItem(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "roadmap item not found"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *RoadmapItemHandler) ListRoadmapItems(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	items, err := h.service.ListRoadmapItems(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": items})
}

type roadmapItemCreateRequest struct {
	Type                domain.RoadmapItemType     `json:"type"`
	Title               string                     `json:"title"`
	Description         string                     `json:"description"`
	BusinessContext     string                     `json:"business_context"`
	TechnicalContext    string                     `json:"technical_context"`
	Priority            domain.RoadmapItemPriority `json:"priority"`
	Status              domain.RoadmapItemStatus   `json:"status"`
	RiskLevel           domain.RiskLevel           `json:"risk_level"`
	BreakingChange      bool                       `json:"breaking_change"`
	RegressionSensitive bool                       `json:"regression_sensitive"`
}

func (h *RoadmapItemHandler) CreateRoadmapItem(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	req := new(roadmapItemCreateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	item := &domain.RoadmapItem{
		ProjectID:           projectID,
		Type:                req.Type,
		Title:               req.Title,
		Description:         req.Description,
		BusinessContext:     req.BusinessContext,
		TechnicalContext:    req.TechnicalContext,
		Priority:            req.Priority,
		Status:              req.Status,
		RiskLevel:           req.RiskLevel,
		BreakingChange:      req.BreakingChange,
		RegressionSensitive: req.RegressionSensitive,
	}
	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	createdItem, err := h.service.CreateRoadmapItem(c.Request().Context(), item, principal.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdItem)
}

type roadmapItemUpdateRequest struct {
	Title            string                   `json:"title"`
	Description      string                   `json:"description"`
	BusinessContext  string                   `json:"business_context"`
	TechnicalContext string                   `json:"technical_context"`
	Status           domain.RoadmapItemStatus `json:"status"`
}

func (h *RoadmapItemHandler) UpdateRoadmapItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}
	req := new(roadmapItemUpdateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	item, err := h.service.UpdateRoadmapItem(c.Request().Context(), id, req.Title, req.Description, req.BusinessContext, req.TechnicalContext, req.Status, principal.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *RoadmapItemHandler) DeleteRoadmapItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if err := h.service.DeleteRoadmapItem(c.Request().Context(), id, principal.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *RoadmapItemHandler) ExportRoadmapItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}

	format := domain.ExportFormat(c.QueryParam("format"))
	if format == "" {
		format = domain.ExportFormatZip
	}

	includeDeps := c.QueryParam("include_dependencies") != "false"
	includeGov := c.QueryParam("include_governance") != "false"

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	pkg, err := h.artifactService.GenerateArtifact(c.Request().Context(), id, format, app.ExportOptions{
		IncludeDependencies: includeDeps,
		IncludeGovernance:   includeGov,
	}, principal.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	data, contentType, err := h.exporter.Export(pkg, format)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if format == domain.ExportFormatZip {
		c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"build-artifact-%s.zip\"", id))
	}

	return c.Blob(http.StatusOK, contentType, data)
}
