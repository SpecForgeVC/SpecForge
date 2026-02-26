import { apiRequest } from "./client";
import type { components } from "./generated/schema";

export type MCPToken = components["schemas"]["MCPToken"];
export type MCPTokenRaw = components["schemas"]["MCPTokenRaw"];
export type MCPStatus = components["schemas"]["MCPStatus"];

export const mcpApi = {
    getStatus: () =>
        apiRequest("/mcp/status", "get"),

    listTokens: (projectId: string) =>
        apiRequest("/mcp/tokens", "get", { params: { project_id: projectId } }),

    generateToken: (projectId: string) =>
        apiRequest("/mcp/tokens", "post", { params: { project_id: projectId } }),

    revokeToken: (tokenId: string) =>
        apiRequest("/mcp/tokens/{id}", "delete", { params: { id: tokenId } }),

    downloadConfig: (ide: string, projectId: string, token?: string) =>
        apiRequest("/mcp/config/download", "get", {
            params: { ide: ide as any, project_id: projectId, token }
        }),
};
