package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type VariableHandler struct {
	service app.VariableService
}

func NewVariableHandler(s app.VariableService) *VariableHandler {
	return &VariableHandler{service: s}
}

func (h *VariableHandler) GetVariable(c echo.Context) error {
	id, err := uuid.Parse(c.Param("variableId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid variable id", err.Error())
	}
	v, err := h.service.GetVariable(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "variable not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, v)
}

func (h *VariableHandler) ListVariables(c echo.Context) error {
	contractID, err := uuid.Parse(c.Param("contractId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid contract id", err.Error())
	}
	vars, err := h.service.ListVariables(c.Request().Context(), contractID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list variables", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, vars)
}

func (h *VariableHandler) ListVariablesByProject(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	vars, err := h.service.ListVariablesByProject(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list variables", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, vars)
}

func (h *VariableHandler) CreateVariable(c echo.Context) error {
	contractID, err := uuid.Parse(c.Param("contractId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid contract id", err.Error())
	}
	var input struct {
		Name            string                 `json:"name"`
		Type            string                 `json:"type"`
		Required        bool                   `json:"required"`
		DefaultValue    string                 `json:"default_value"`
		Description     string                 `json:"description"`
		ValidationRules map[string]interface{} `json:"validation_rules"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	v, err := h.service.CreateVariable(c.Request().Context(), contractID, input.Name, input.Type, input.Required, input.DefaultValue, input.Description, input.ValidationRules, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create variable", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, v)
}

func (h *VariableHandler) CreateVariableByProject(c echo.Context) error {
	var input struct {
		ContractID      uuid.UUID              `json:"contract_id"`
		Name            string                 `json:"name"`
		Type            string                 `json:"type"`
		Required        bool                   `json:"required"`
		DefaultValue    string                 `json:"default_value"`
		Description     string                 `json:"description"`
		ValidationRules map[string]interface{} `json:"validation_rules"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	if input.ContractID == uuid.Nil {
		return ErrorResponse(c, http.StatusBadRequest, "MISSING_FIELD", "contract_id is required", "")
	}
	userID := GetUserID(c)
	v, err := h.service.CreateVariable(c.Request().Context(), input.ContractID, input.Name, input.Type, input.Required, input.DefaultValue, input.Description, input.ValidationRules, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create variable", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, v)
}

func (h *VariableHandler) UpdateVariable(c echo.Context) error {
	id, err := uuid.Parse(c.Param("variableId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid variable id", err.Error())
	}
	var input struct {
		Name            string                 `json:"name"`
		Type            string                 `json:"type"`
		Required        bool                   `json:"required"`
		DefaultValue    string                 `json:"default_value"`
		Description     string                 `json:"description"`
		ValidationRules map[string]interface{} `json:"validation_rules"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	v, err := h.service.UpdateVariable(c.Request().Context(), id, input.Name, input.Type, input.Required, input.DefaultValue, input.Description, input.ValidationRules, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update variable", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, v)
}

func (h *VariableHandler) DeleteVariable(c echo.Context) error {
	id, err := uuid.Parse(c.Param("variableId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid variable id", err.Error())
	}
	userID := GetUserID(c)
	if err := h.service.DeleteVariable(c.Request().Context(), id, userID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete variable", err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
