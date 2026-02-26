package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type FeatureIntelligenceHandler struct {
	service app.FeatureIntelligenceService
}

func NewFeatureIntelligenceHandler(service app.FeatureIntelligenceService) *FeatureIntelligenceHandler {
	return &FeatureIntelligenceHandler{service: service}
}

func (h *FeatureIntelligenceHandler) GetFeatureIntelligence(c echo.Context) error {
	id, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid roadmap item id", err.Error())
	}

	// Try to get existing score
	fi, err := h.service.GetFeatureScore(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get feature intelligence", err.Error())
	}

	if fi == nil {
		// Calculate if not found
		fi, err = h.service.CalculateFeatureScore(c.Request().Context(), id)
		if err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to calculate feature intelligence", err.Error())
		}
	}

	return SuccessResponse(c, http.StatusOK, fi)
}
