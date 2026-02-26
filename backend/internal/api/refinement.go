package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
)

type RefinementHandler struct {
	service app.RefinementService
}

func NewRefinementHandler(service app.RefinementService) *RefinementHandler {
	return &RefinementHandler{service: service}
}

type StartSessionRequest struct {
	ArtifactType  string         `json:"artifact_type"`
	TargetType    string         `json:"target_type"`
	Prompt        string         `json:"prompt"`
	ContextData   map[string]any `json:"context_data"`
	MaxIterations int            `json:"max_iterations"`
}

func (h *RefinementHandler) StartSession(c echo.Context) error {
	var req StartSessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	session, err := h.service.StartSession(c.Request().Context(), req.ArtifactType, req.TargetType, req.Prompt, req.ContextData, req.MaxIterations)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, session)
}

func (h *RefinementHandler) StreamEvents(c echo.Context) error {
	id, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session id"})
	}

	eventCh, err := h.service.GetSessionEvents(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	// Flush initial headers
	c.Response().Flush()

	timeout := time.After(5 * time.Minute) // Safety timeout

	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				// Channel closed by orchestrator
				fmt.Fprintf(c.Response(), "event: done\ndata: {}\n\n")
				c.Response().Flush()
				return nil
			}
			data, _ := json.Marshal(event)
			fmt.Fprintf(c.Response(), "data: %s\n\n", string(data))
			c.Response().Flush()
		case <-c.Request().Context().Done():
			return nil
		case <-timeout:
			return nil
		}
	}
}

func (h *RefinementHandler) ApproveSession(c echo.Context) error {
	id, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session id"})
	}

	if err := h.service.ApproveSession(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "approved"})
}
