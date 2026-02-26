package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type AuditLogHandler struct {
	service app.AuditLogService
}

func NewAuditLogHandler(service app.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

func (h *AuditLogHandler) GetEntityLogs(c echo.Context) error {
	entityType := c.Param("entityType")
	entityID, err := uuid.Parse(c.Param("entityId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid entity id"})
	}

	logs, err := h.service.GetEntityLogs(c.Request().Context(), entityType, entityID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": logs})
}

func (h *AuditLogHandler) GetRoadmapItemActivity(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}

	logs, err := h.service.GetEntityLogs(c.Request().Context(), "ROADMAP_ITEM", roadmapItemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": true, "data": logs})
}
