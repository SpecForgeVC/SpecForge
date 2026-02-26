package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type ValidationRuleHandler struct {
	service app.ValidationRuleService
}

func NewValidationRuleHandler(s app.ValidationRuleService) *ValidationRuleHandler {
	return &ValidationRuleHandler{service: s}
}

func (h *ValidationRuleHandler) GetValidationRule(c echo.Context) error {
	id, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid rule id", err.Error())
	}
	rule, err := h.service.GetValidationRule(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "validation rule not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, rule)
}

func (h *ValidationRuleHandler) ListValidationRules(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	rules, err := h.service.ListValidationRules(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list validation rules", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, rules)
}

func (h *ValidationRuleHandler) CreateValidationRule(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	var input struct {
		Name        string                 `json:"name"`
		RuleType    string                 `json:"rule_type"`
		RuleConfig  map[string]interface{} `json:"rule_config"`
		Description string                 `json:"description"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	rule, err := h.service.CreateValidationRule(c.Request().Context(), projectID, input.Name, input.RuleType, input.RuleConfig, input.Description, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create validation rule", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, rule)
}

func (h *ValidationRuleHandler) UpdateValidationRule(c echo.Context) error {
	id, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid rule id", err.Error())
	}
	var input struct {
		Name        string                 `json:"name"`
		RuleType    string                 `json:"rule_type"`
		RuleConfig  map[string]interface{} `json:"rule_config"`
		Description string                 `json:"description"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	rule, err := h.service.UpdateValidationRule(c.Request().Context(), id, input.Name, input.RuleType, input.RuleConfig, input.Description, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update validation rule", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, rule)
}

func (h *ValidationRuleHandler) DeleteValidationRule(c echo.Context) error {
	id, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid rule id", err.Error())
	}
	userID := GetUserID(c)
	if err := h.service.DeleteValidationRule(c.Request().Context(), id, userID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete validation rule", err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
