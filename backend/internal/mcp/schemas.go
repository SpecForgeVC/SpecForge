package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/scott/specforge/internal/domain"
)

// ToolInputCreateSnapshot defines the input for the create_snapshot tool
type ToolInputCreateSnapshot struct {
	ProjectID     string `json:"project_id"`
	RoadmapItemID string `json:"roadmap_item_id"`
	Mode          string `json:"mode"`
}

// ToolOutputCreateSnapshot defines the output for the create_snapshot tool
type ToolOutputCreateSnapshot struct {
	SnapshotID             string      `json:"snapshot_id"`
	ExtractionRequirements interface{} `json:"extraction_requirements"`
	RequiredSchema         interface{} `json:"required_schema"`
	StrictNextStep         string      `json:"strict_next_step"`
}

// ToolInputPostSnapshot defines the input for the post_snapshot tool
type ToolInputPostSnapshot struct {
	SnapshotID          string                     `json:"snapshot_id"`
	EnvironmentSnapshot domain.EnvironmentSnapshot `json:"environment_snapshot"`
}

// ToolOutputPostSnapshot defines the output for the post_snapshot tool
type ToolOutputPostSnapshot struct {
	AnalysisResults  interface{}   `json:"analysis_results"`
	Scores           domain.Scores `json:"scores"`
	Verdict          string        `json:"verdict"` // approved, blocked, requires_alignment
	RequiredNextTool string        `json:"required_next_tool"`
}

// JSON-RPC models
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{}   `json:"id"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *JSONRPCError) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}
