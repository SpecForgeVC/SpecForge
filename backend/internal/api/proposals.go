package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
	mw "github.com/scott/specforge/internal/transport/middleware"
)

type AiProposalHandler struct {
	service app.AiProposalService
}

func NewAiProposalHandler(service app.AiProposalService) *AiProposalHandler {
	return &AiProposalHandler{service: service}
}

func (h *AiProposalHandler) GetProposal(c echo.Context) error {
	id, err := uuid.Parse(c.Param("proposalId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proposal id"})
	}
	proposal, err := h.service.GetProposal(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "proposal not found"})
	}
	return c.JSON(http.StatusOK, proposal)
}

func (h *AiProposalHandler) ListProposals(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid project id"})
	}
	proposals, err := h.service.ListProposals(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": proposals})
}

type proposalCreateRequest struct {
	RoadmapItemID   uuid.UUID              `json:"roadmap_item_id"`
	ProposalType    domain.ProposalType    `json:"proposal_type"`
	Diff            map[string]interface{} `json:"diff"`
	Reasoning       string                 `json:"reasoning"`
	ConfidenceScore float64                `json:"confidence_score"`
}

func (h *AiProposalHandler) CreateProposal(c echo.Context) error {
	req := new(proposalCreateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	proposal, err := h.service.CreateProposal(c.Request().Context(), req.RoadmapItemID, req.ProposalType, req.Diff, req.Reasoning, req.ConfidenceScore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, proposal)
}

func (h *AiProposalHandler) ApproveProposal(c echo.Context) error {
	id, err := uuid.Parse(c.Param("proposalId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proposal id"})
	}

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if err := h.service.ApproveProposal(c.Request().Context(), id, principal.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *AiProposalHandler) RejectProposal(c echo.Context) error {
	id, err := uuid.Parse(c.Param("proposalId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proposal id"})
	}

	principal, ok := mw.PrincipalFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if err := h.service.RejectProposal(c.Request().Context(), id, principal.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
