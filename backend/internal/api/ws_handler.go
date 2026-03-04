package api

import (
	"net/http"

	"github.com/SpecForgeVC/SpecForge/internal/app"
	"github.com/SpecForgeVC/SpecForge/internal/infra/auth"
	"github.com/SpecForgeVC/SpecForge/internal/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{}

type WSHandler struct {
	service        app.NotificationService
	validator      *auth.JWTValidator
	allowedOrigins []string
}

func NewWSHandler(s app.NotificationService, v *auth.JWTValidator, allowedOrigins []string) *WSHandler {
	return &WSHandler{
		service:        s,
		validator:      v,
		allowedOrigins: allowedOrigins,
	}
}

func (h *WSHandler) Connect(c echo.Context) error {
	// 1. Extract token from query param
	tokenStr := c.QueryParam("token")
	if tokenStr == "" {
		logger.Warn("WebSocket connection attempt without token")
		return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
	}

	// 2. Validate token
	principal, err := h.validator.Validate(tokenStr)
	if err != nil {
		logger.Warn("WebSocket connection attempt with invalid token", zap.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}

	// 3. Upgrade HTTP connection to WebSocket
	upgrader.CheckOrigin = func(r *http.Request) bool {
		if len(h.allowedOrigins) == 0 {
			return true // Fallback to allow all if not configured (dev)
		}
		origin := r.Header.Get("Origin")
		for _, o := range h.allowedOrigins {
			if o == origin {
				return true
			}
		}
		return false
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.Error("Failed to upgrade websocket", zap.Error(err))
		return err
	}

	userID := principal.UserID
	h.service.Register(ws, userID)
	logger.Info("Client connected", zap.String("userID", userID.String()))

	defer func() {
		h.service.Unregister(ws)
		logger.Info("Client disconnected", zap.String("userID", userID.String()))
	}()

	// Read loop to keep connection open and handle close
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}
