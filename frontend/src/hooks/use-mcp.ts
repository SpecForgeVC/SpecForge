import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { mcpApi, type MCPTokenRaw } from "@/api/mcp";

export const mcpKeys = {
    all: ["mcp"] as const,
    tokens: (projectId: string) => [...mcpKeys.all, "tokens", projectId] as const,
};

export function useMCPTokens(projectId?: string) {
    return useQuery({
        queryKey: mcpKeys.tokens(projectId || ""),
        queryFn: () => mcpApi.listTokens(projectId!),
        enabled: !!projectId,
    });
}

export function useGenerateMCPToken(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation<MCPTokenRaw, Error, void>({
        mutationFn: () => mcpApi.generateToken(projectId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: mcpKeys.tokens(projectId) });
        },
    });
}

export function useRevokeMCPToken(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation<void, Error, string>({
        mutationFn: (tokenId: string) => mcpApi.revokeToken(tokenId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: mcpKeys.tokens(projectId) });
        },
    });
}
