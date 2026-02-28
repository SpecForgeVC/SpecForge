package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/ui_roadmap"
)

type UIRoadmapHandler struct {
	service ui_roadmap.Service
}

func NewUIRoadmapHandler(service ui_roadmap.Service) *UIRoadmapHandler {
	return &UIRoadmapHandler{service: service}
}

func (h *UIRoadmapHandler) ListItems(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	items, err := h.service.ListItems(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *UIRoadmapHandler) GetItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	item, err := h.service.GetItem(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "item not found"})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *UIRoadmapHandler) SaveItem(c echo.Context) error {
	var item ui_roadmap.UIRoadmapItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Ensure project ID from URL if provided
	if pID := c.Param("projectId"); pID != "" {
		projectID, err := uuid.Parse(pID)
		if err == nil {
			item.ProjectID = projectID
		}
	}

	if err := h.service.SaveItem(c.Request().Context(), &item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, item)
}

func (h *UIRoadmapHandler) ExportItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	bundle, err := h.service.Export(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, bundle)
}

func (h *UIRoadmapHandler) DeleteItem(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.service.DeleteItem(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *UIRoadmapHandler) SyncFigma(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var payload ui_roadmap.FigmaSyncPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.service.SyncFigma(c.Request().Context(), id, payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *UIRoadmapHandler) GetPluginAssets(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	assets, err := h.service.GetPluginAssets(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, assets)
}

func (h *UIRoadmapHandler) RecommendTree(c echo.Context) error {
	var item ui_roadmap.UIRoadmapItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	tree, err := h.service.RecommendComponentTree(c.Request().Context(), &item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"component_tree": tree})
}

func (h *UIRoadmapHandler) RecommendStateMachine(c echo.Context) error {
	var item ui_roadmap.UIRoadmapItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sm, err := h.service.RecommendStateMachine(c.Request().Context(), &item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, sm)
}

func (h *UIRoadmapHandler) RecommendFix(c echo.Context) error {
	var req struct {
		Item   ui_roadmap.UIRoadmapItem `json:"item"`
		Issues []ui_roadmap.DriftIssue  `json:"issues"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	fixedItem, err := h.service.RecommendFix(c.Request().Context(), &req.Item, req.Issues)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, fixedItem)
}

func (h *UIRoadmapHandler) CheckCompliance(c echo.Context) error {
	var item ui_roadmap.UIRoadmapItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	issues, err := h.service.CheckCompliance(c.Request().Context(), &item)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"issues": issues})
}

func (h *UIRoadmapHandler) RecommendAPIContracts(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}

	rec, err := h.service.RecommendAPIContracts(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rec)
}
