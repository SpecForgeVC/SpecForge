package api

import (
	"github.com/labstack/echo/v4"
)

// Response is the standardized JSON response envelope.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// Error details.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse returns a successful response without metadata.
func SuccessResponse(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessResponseWithMeta returns a successful response with metadata.
func SuccessResponseWithMeta(c echo.Context, status int, data interface{}, meta interface{}) error {
	return c.JSON(status, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse returns an error response.
func ErrorResponse(c echo.Context, status int, code, message, details string) error {
	return c.JSON(status, Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
