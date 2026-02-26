package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Adapt converts a standard net/http middleware to Echo middleware.
func Adapt(m func(http.Handler) http.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var err error
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				err = next(c)
			})).ServeHTTP(c.Response().Writer, c.Request())
			return err
		}
	}
}
