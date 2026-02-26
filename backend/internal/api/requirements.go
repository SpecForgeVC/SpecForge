package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type RequirementHandler struct {
	service app.RequirementService
}

func NewRequirementHandler(s app.RequirementService) *RequirementHandler {
	return &RequirementHandler{service: s}
}

func (h *RequirementHandler) GetRequirement(c echo.Context) error {
	id, err := uuid.Parse(c.Param("requirementId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid requirement id", err.Error())
	}
	req, err := h.service.GetRequirement(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "requirement not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, req)
}

func (h *RequirementHandler) ListRequirements(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid roadmap item id", err.Error())
	}
	reqs, err := h.service.ListRequirements(c.Request().Context(), roadmapItemID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list requirements", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, reqs)
}

func (h *RequirementHandler) CreateRequirement(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid roadmap item id", err.Error())
	}
	var input struct {
		Title              string `json:"title"`
		Description        string `json:"description"`
		Testable           bool   `json:"testable"`
		AcceptanceCriteria string `json:"acceptance_criteria"`
		OrderIndex         int    `json:"order_index"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	req, err := h.service.CreateRequirement(c.Request().Context(), roadmapItemID, input.Title, input.Description, input.Testable, input.AcceptanceCriteria, input.OrderIndex, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create requirement", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, req)
}

func (h *RequirementHandler) UpdateRequirement(c echo.Context) error {
	id, err := uuid.Parse(c.Param("requirementId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid requirement id", err.Error())
	}
	var input struct {
		Title              string `json:"title"`
		Description        string `json:"description"`
		Testable           bool   `json:"testable"`
		AcceptanceCriteria string `json:"acceptance_criteria"`
		OrderIndex         int    `json:"order_index"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	req, err := h.service.UpdateRequirement(c.Request().Context(), id, input.Title, input.Description, input.Testable, input.AcceptanceCriteria, input.OrderIndex, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update requirement", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, req)
}

func (h *RequirementHandler) DeleteRequirement(c echo.Context) error {
	id, err := uuid.Parse(c.Param("requirementId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid requirement id", err.Error())
	}
	userID := GetUserID(c)
	if err := h.service.DeleteRequirement(c.Request().Context(), id, userID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete requirement", err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
