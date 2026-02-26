package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Router handles incoming MCP JSON-RPC requests
type Router struct {
	handlers map[string]HandlerFunc
	server   *Server
}

// HandlerFunc is the signature for MCP tool handlers
type HandlerFunc func(ctx context.Context, params json.RawMessage) (interface{}, error)

func NewRouter(server *Server) *Router {
	r := &Router{
		handlers: make(map[string]HandlerFunc),
		server:   server,
	}
	r.registerHandlers()
	return r
}

func (r *Router) registerHandlers() {
	r.handlers["tools/list"] = r.handleListTools
	r.handlers["tools/call"] = r.handleCallTool
	r.handlers["initialize"] = r.handleInitialize
	r.handlers["notifications/initialized"] = r.handleInitialized
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var rpcReq JSONRPCRequest
	if err := json.NewDecoder(req.Body).Decode(&rpcReq); err != nil {
		r.sendError(w, nil, -32700, "Parse error", nil)
		return
	}

	// AUTHENTICATION
	if r.server.config.AuthRequired {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			r.sendError(w, rpcReq.ID, -32000, "Authentication required", nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			r.sendError(w, rpcReq.ID, -32000, "Invalid authorization format", nil)
			return
		}

		token := parts[1]

		// Fallback to static token for legacy support if configured
		if token != r.server.config.Token {
			// Validate against dynamic token system
			if r.server.config.TokenService != nil {
				_, err := r.server.config.TokenService.ValidateToken(req.Context(), token)
				if err != nil {
					r.sendError(w, rpcReq.ID, -32000, "Unauthorized", nil)
					return
				}
			} else {
				r.sendError(w, rpcReq.ID, -32000, "Unauthorized", nil)
				return
			}
		}
	}

	result, err := r.Route(req.Context(), rpcReq)
	if err != nil {
		if rpcErr, ok := err.(*JSONRPCError); ok {
			r.sendError(w, rpcReq.ID, rpcErr.Code, rpcErr.Message, rpcErr.Data)
		} else {
			r.sendError(w, rpcReq.ID, -32603, "Internal error", err.Error())
		}
		return
	}

	r.sendResult(w, rpcReq.ID, result)
}

func (r *Router) Route(ctx context.Context, req JSONRPCRequest) (interface{}, error) {
	handler, ok := r.handlers[req.Method]
	if !ok {
		return nil, &JSONRPCError{Code: -32601, Message: "Method not found"}
	}
	return handler(ctx, req.Params)
}

func (r *Router) sendResult(w http.ResponseWriter, id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (r *Router) sendError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Handler implementations for core MCP lifecycle
func (r *Router) handleInitialize(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "SpecForge Reality Anchor Engine",
			"version": "1.0.0",
		},
	}, nil
}

func (r *Router) handleListTools(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{
		"tools": GetToolDefinitions(),
	}, nil
}

func (r *Router) handleInitialized(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return map[string]interface{}{}, nil
}

func (r *Router) handleCallTool(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var callParams struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, &JSONRPCError{Code: -32602, Message: "Invalid params"}
	}

	// Tool-specific routing
	switch callParams.Name {
	case "create_snapshot":
		return r.server.handlers.CreateSnapshot(ctx, callParams.Arguments)
	case "post_snapshot":
		return r.server.handlers.PostSnapshot(ctx, callParams.Arguments)
	case "get_snapshot_status":
		return r.server.handlers.GetSnapshotStatus(ctx, callParams.Arguments)
	case "list_active_snapshots":
		return r.server.handlers.ListActiveSnapshots(ctx, callParams.Arguments)
	case "init_project_import":
		return r.server.handlers.InitProjectImport(ctx, callParams.Arguments)
	case "submit_project_snapshot":
		return r.server.handlers.SubmitProjectSnapshot(ctx, callParams.Arguments)
	case "get_import_alignment_rules":
		return r.server.handlers.GetImportAlignmentRules(ctx, callParams.Arguments)
	case "submit_post_import_snapshot":
		return r.server.handlers.SubmitPostImportSnapshot(ctx, callParams.Arguments)
	case "finalize_project_import":
		return r.server.handlers.FinalizeProjectImport(ctx, callParams.Arguments)
	case "help":
		return r.server.handlers.Help(ctx, callParams.Arguments)
	default:
		return nil, &JSONRPCError{Code: -32601, Message: fmt.Sprintf("Tool not found: %s", callParams.Name)}
	}
}
