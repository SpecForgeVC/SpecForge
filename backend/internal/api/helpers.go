package api

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/transport/middleware"
)

// GetUserID extracts the user ID from the Principal in the request context.
func GetUserID(c echo.Context) uuid.UUID {
	p, ok := middleware.PrincipalFromContext(c.Request().Context())
	if !ok {
		return uuid.Nil
	}
	return p.UserID
}
