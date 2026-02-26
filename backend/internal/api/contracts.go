package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/scott/specforge/internal/app"
	"github.com/scott/specforge/internal/domain"
)

type ContractHandler struct {
	service app.ContractService
}

func NewContractHandler(service app.ContractService) *ContractHandler {
	return &ContractHandler{service: service}
}

func (h *ContractHandler) GetContract(c echo.Context) error {
	id, err := uuid.Parse(c.Param("contractId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid contract id", err.Error())
	}
	contract, err := h.service.GetContract(c.Request().Context(), id)
	if err != nil {
		return ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "contract not found", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, contract)
}

func (h *ContractHandler) ListContracts(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid roadmap item id", err.Error())
	}
	contracts, err := h.service.ListContracts(c.Request().Context(), roadmapItemID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list contracts", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, contracts)
}

func (h *ContractHandler) ListContractsByProject(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid project id", err.Error())
	}
	contracts, err := h.service.ListContractsByProject(c.Request().Context(), projectID)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list contracts", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, contracts)
}

type contractCreateRequest struct {
	ContractType domain.ContractType    `json:"contract_type"`
	Version      string                 `json:"version"`
	InputSchema  map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema"`
	ErrorSchema  map[string]interface{} `json:"error_schema"`
}

func (h *ContractHandler) CreateContract(c echo.Context) error {
	roadmapItemID, err := uuid.Parse(c.Param("roadmapItemId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid roadmap item id", err.Error())
	}
	req := new(contractCreateRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	contract, err := h.service.CreateContract(c.Request().Context(), roadmapItemID, req.ContractType, req.Version, req.InputSchema, req.OutputSchema, req.ErrorSchema)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create contract", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, contract)
}

type contractCreateByProjectRequest struct {
	RoadmapItemID uuid.UUID              `json:"roadmap_item_id"`
	ContractType  domain.ContractType    `json:"contract_type"`
	Version       string                 `json:"version"`
	InputSchema   map[string]interface{} `json:"input_schema"`
	OutputSchema  map[string]interface{} `json:"output_schema"`
	ErrorSchema   map[string]interface{} `json:"error_schema"`
}

func (h *ContractHandler) CreateContractByProject(c echo.Context) error {
	// projectId is in URL but not strictly needed if roadmap_item_id is provided,
	// but we can use it for validation if needed.
	req := new(contractCreateByProjectRequest)
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	if req.RoadmapItemID == uuid.Nil {
		return ErrorResponse(c, http.StatusBadRequest, "MISSING_FIELD", "roadmap_item_id is required", "")
	}
	contract, err := h.service.CreateContract(c.Request().Context(), req.RoadmapItemID, req.ContractType, req.Version, req.InputSchema, req.OutputSchema, req.ErrorSchema)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create contract", err.Error())
	}
	return SuccessResponse(c, http.StatusCreated, contract)
}

func (h *ContractHandler) UpdateContract(c echo.Context) error {
	id, err := uuid.Parse(c.Param("contractId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid contract id", err.Error())
	}
	req := new(contractCreateRequest) // Reuse create request struct as fields are same
	if err := c.Bind(req); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_BODY", "failed to bind request", err.Error())
	}
	contract, err := h.service.UpdateContract(c.Request().Context(), id, req.ContractType, req.Version, req.InputSchema, req.OutputSchema, req.ErrorSchema)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update contract", err.Error())
	}
	return SuccessResponse(c, http.StatusOK, contract)
}

func (h *ContractHandler) DeleteContract(c echo.Context) error {
	id, err := uuid.Parse(c.Param("contractId"))
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "invalid contract id", err.Error())
	}
	if err := h.service.DeleteContract(c.Request().Context(), id); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete contract", err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
