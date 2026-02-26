package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type WebhookHandler struct {
	service app.WebhookService
}

func NewWebhookHandler(s app.WebhookService) *WebhookHandler {
	return &WebhookHandler{service: s}
}

func (h *WebhookHandler) GetWebhook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("webhookId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid webhook id", err.Error())
	}
	w, err := h.service.GetWebhook(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "webhook not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, w)
}

func (h *WebhookHandler) ListWebhooks(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	webhooks, err := h.service.ListWebhooks(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list webhooks", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, webhooks)
}

func (h *WebhookHandler) CreateWebhook(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	var input struct {
		URL    string   `json:"url"`
		Events []string `json:"events"`
		Secret string   `json:"secret"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	w, err := h.service.CreateWebhook(c.Request().Context(), projectID, input.URL, input.Events, input.Secret, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create webhook", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, w)
}

func (h *WebhookHandler) UpdateWebhook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("webhookId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid webhook id", err.Error())
	}
	var input struct {
		URL    string   `json:"url"`
		Events []string `json:"events"`
		Secret string   `json:"secret"`
		Active bool     `json:"active"`
	}
	if err := c.Bind(&input); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	userID := GetUserID(c)
	w, err := h.service.UpdateWebhook(c.Request().Context(), id, input.URL, input.Events, input.Secret, input.Active, userID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update webhook", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, w)
}

func (h *WebhookHandler) DeleteWebhook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("webhookId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid webhook id", err.Error())
	}
	userID := GetUserID(c)
	if err := h.service.DeleteWebhook(c.Request().Context(), id, userID); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete webhook", err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
