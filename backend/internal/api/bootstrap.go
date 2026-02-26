package api

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type BootstrapHandler struct {
	service app.BootstrapService
}

func NewBootstrapHandler(service app.BootstrapService) *BootstrapHandler {
	return &BootstrapHandler{service: service}
}

func (h *BootstrapHandler) GeneratePrompt(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}

	prompt, projectName, err := h.service.GeneratePrompt(c.Request().Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "project not found", err.Error())
		}
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate prompt", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, map[string]interface{}{
		"prompt":       prompt,
		"project_name": projectName,
		"project_id":   projectID,
	})
}

func (h *BootstrapHandler) IngestBootstrap(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}

	var payload domain.BootstrapPayload
	if err := c.Bind(&payload); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "failed to parse bootstrap payload", err.Error())
	}

	snapshot, warnings, err := h.service.IngestBootstrap(c.Request().Context(), projectID, payload)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "project not found", err.Error())
		}
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to ingest bootstrap", err.Error())
	}

	return SuccessResponse(c, http.StatusCreated, map[string]interface{}{
		"snapshot": snapshot,
		"scores": map[string]interface{}{
			"architecture_score": snapshot.ArchitectureScore,
			"contract_density":   snapshot.ContractDensity,
			"risk_score":         snapshot.RiskScore,
			"alignment_score":    snapshot.AlignmentScore,
		},
		"confidence": snapshot.ConfidenceJSON,
		"warnings":   warnings,
	})
}

func (h *BootstrapHandler) ListSnapshots(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}

	snapshots, err := h.service.ListSnapshots(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list snapshots", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, snapshots)
}

func (h *BootstrapHandler) GetSnapshot(c echo.Context) error {
	snapshotID, err := uuid.Parse(c.Param("snapshotId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid snapshot id", err.Error())
	}

	snapshot, err := h.service.GetSnapshot(c.Request().Context(), snapshotID)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "snapshot not found", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, snapshot)
}

func (h *BootstrapHandler) GetLatestSnapshot(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}

	snapshot, err := h.service.GetLatestSnapshot(c.Request().Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "project not found", err.Error())
		}
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "no snapshots found", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, snapshot)
}

func (h *BootstrapHandler) GetLatestImportSession(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}

	session, err := h.service.GetLatestImportSession(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to fetch session", err.Error())
	}

	if session == nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "no import session found", "no import session found")
	}

	return SuccessResponse(c, http.StatusOK, session)
}

type diffRequest struct {
	FromSnapshotID *uuid.UUID `json:"from_snapshot_id"`
	ToSnapshotID   *uuid.UUID `json:"to_snapshot_id"`
}

func (h *BootstrapHandler) DiffSnapshots(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}

	var req diffRequest
	if err := c.Bind(&req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "failed to parse diff request", err.Error())
	}

	diff, err := h.service.DiffSnapshots(c.Request().Context(), projectID, req.FromSnapshotID, req.ToSnapshotID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to diff snapshots", err.Error())
	}

	return SuccessResponse(c, http.StatusOK, diff)
}
