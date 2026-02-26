package api

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Pagination stores pagination parameters.
type Pagination struct {
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
	Offset   int `json:"offset"`
}

// GetPagination extracts pagination parameters from the query strings.
func GetPagination(c echo.Context) Pagination {
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	page, _ := strconv.Atoi(c.QueryParam("page"))

	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}

	return Pagination{
		PageSize: pageSize,
		Page:     page,
		Offset:   (page - 1) * pageSize,
	}
}
