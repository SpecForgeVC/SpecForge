package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type LLMSettingsHandler struct {
	service app.LLMService
}

func NewLLMSettingsHandler(service app.LLMService) *LLMSettingsHandler {
	return &LLMSettingsHandler{service: service}
}

func (h *LLMSettingsHandler) GetConfig(c echo.Context) error {
	config, err := h.service.GetActiveConfig(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if config == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"data": nil})
	}
	// Mask API Key before returning
	config.APIKey = "********"
	return c.JSON(http.StatusOK, config)
}

func (h *LLMSettingsHandler) UpdateConfig(c echo.Context) error {
	var req domain.LLMConfiguration
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := h.service.UpdateConfig(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, req)
}

func (h *LLMSettingsHandler) TestConnection(c echo.Context) error {
	var req domain.LLMConfiguration
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := h.service.TestConfiguration(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Connection failed: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Connection successful"})
}

func (h *LLMSettingsHandler) Warmup(c echo.Context) error {
	ctx := c.Request().Context()
	client, err := h.service.GetClient(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	chunks := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(chunks)
		defer close(errCh)
		if err := client.StreamGenerate(ctx, "Hello, can you confirm you are working?", chunks); err != nil {
			errCh <- err
		}
	}()

	ticker := time.NewTicker(100 * time.Millisecond) // Keep-alive ticker if needed
	defer ticker.Stop()

	for {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				fmt.Fprintf(c.Response(), "event: done\ndata: {}\n\n")
				c.Response().Flush()
				return nil
			}
			fmt.Fprintf(c.Response(), "data: %s\n\n", chunk)
			c.Response().Flush()
		case err := <-errCh:
			if err != nil {
				fmt.Fprintf(c.Response(), "event: error\ndata: %s\n\n", err.Error())
				c.Response().Flush()
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (h *LLMSettingsHandler) ListModels(c echo.Context) error {
	var req domain.LLMConfiguration
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	models, err := h.service.ListModels(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Failed to list models: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"models": models})
}
