package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/domain"
	"github.com/scott/specforge/internal/drift"
)

type DriftHandler struct {
	service drift.DriftService
}

func NewDriftHandler(s drift.DriftService) *DriftHandler {
	return &DriftHandler{service: s}
}

func (h *DriftHandler) RunDriftCheck(c echo.Context) error {
	contractID, err := uuid.Parse(c.Param("contractId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contractId"})
	}

	var req struct {
		AgainstVersion uuid.UUID `json:"against_version"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	report, err := h.service.RunDriftCheck(c.Request().Context(), contractID, req.AgainstVersion)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, report)
}

func (h *DriftHandler) GetDriftHistory(c echo.Context) error {
	history, err := h.service.GetDriftHistory(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, history)
}

func (h *DriftHandler) GenerateDriftFixes(c echo.Context) error {
	var req struct {
		DriftReport struct {
			DriftDetected   bool `json:"drift_detected"`
			BreakingChanges []struct {
				Field string `json:"field"`
				Issue string `json:"issue"`
			} `json:"breaking_changes"`
			RiskScore float64 `json:"risk_score"`
		} `json:"drift_report"`
		RoadmapItemID uuid.UUID `json:"roadmap_item_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// Convert to domain types
	report := &domain.DriftReport{
		DriftDetected:   req.DriftReport.DriftDetected,
		RiskScore:       req.DriftReport.RiskScore,
		BreakingChanges: make([]domain.BreakingChange, len(req.DriftReport.BreakingChanges)),
	}
	for i, bc := range req.DriftReport.BreakingChanges {
		report.BreakingChanges[i] = domain.BreakingChange{Field: bc.Field, Issue: bc.Issue}
	}

	fixes, err := h.service.GenerateDriftFixes(c.Request().Context(), report, req.RoadmapItemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"fixes": fixes})
}
