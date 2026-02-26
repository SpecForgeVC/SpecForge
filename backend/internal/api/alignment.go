package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type AlignmentHandler struct {
	service app.AlignmentService
}

func NewAlignmentHandler(service app.AlignmentService) *AlignmentHandler {
	return &AlignmentHandler{service: service}
}

func (h *AlignmentHandler) GetAlignmentReport(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}

	report, err := h.service.GetAlignmentReport(c.Request().Context(), projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if report == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No alignment report found"})
	}

	return c.JSON(http.StatusOK, report)
}

func (h *AlignmentHandler) TriggerAlignmentCheck(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}

	report, err := h.service.TriggerAlignmentCheck(c.Request().Context(), projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, report)
}
