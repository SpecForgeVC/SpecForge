package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xeipuuv/gojsonschema"
)

func SchemaValidationMiddleware(schema map[string]interface{}) echo.MiddlewareFunc {
	schemaLoader := gojsonschema.NewGoLoader(schema)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Read body
			bodyBytes, _ := io.ReadAll(c.Request().Body)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			documentLoader := gojsonschema.NewBytesLoader(bodyBytes)
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			if !result.Valid() {
				var errors []string
				for _, desc := range result.Errors() {
					errors = append(errors, desc.String())
				}
				return c.JSON(http.StatusBadRequest, map[string]interface{}{"errors": errors})
			}

			return next(c)
		}
	}
}
