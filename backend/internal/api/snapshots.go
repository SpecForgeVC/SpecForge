package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	mw "github.com/scott/specforge/internal/transport/middleware"
)

type SnapshotHandler struct {
	service app.SnapshotService
}

func NewSnapshotHandler(service app.SnapshotService) *SnapshotHandler {
	return &SnapshotHandler{service: service}
}

func (h *SnapshotHandler) GetSnapshot(c echo.Context) error {
	id, err := uuid.Parse(c.Param("snapshotId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid snapshot id", err.Error())
	}
	snap, err := h.service.GetSnapshot(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "snapshot not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, snap)
}

func (h *SnapshotHandler) ListSnapshots(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid roadmap item id", err.Error())
	}
	snapshots, err := h.service.ListSnapshots(c.Request().Context(), roadmapItemID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list snapshots", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, snapshots)
}

func (h *SnapshotHandler) ListSnapshotsByProject(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	snapshots, err := h.service.ListSnapshotsByProject(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list snapshots", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, snapshots)
}

type snapshotCreateRequest struct {
	SnapshotData map[string]interface{} `json:"snapshot_data"`
}

func (h *SnapshotHandler) CreateSnapshot(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid roadmap item id"})
	}
	req := new(snapshotCreateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	snap, err := h.service.CreateSnapshot(c.Request().Context(), roadmapItemID, req.SnapshotData, principal.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, snap)
}
