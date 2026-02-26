package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type VariableLineageHandler struct {
	service app.VariableLineageService
}

func NewVariableLineageHandler(service app.VariableLineageService) *VariableLineageHandler {
	return &VariableLineageHandler{service: service}
}

func (h *VariableLineageHandler) GetLineageEvents(c echo.Context) error {
	id, err := uuid.Parse(c.Param("variableId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid variable id", err.Error())
	}

	events, err := h.service.GetLineageEvents(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get lineage events", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, events)
}

func (h *VariableLineageHandler) GetLineageGraph(c echo.Context) error {
	id, err := uuid.Parse(c.Param("variableId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid variable id", err.Error())
	}

	graph, err := h.service.GetLineageGraph(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get lineage graph", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, graph)
}
