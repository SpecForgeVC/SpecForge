import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useWorkspaces() {
    return useQuery({
        queryKey: ["workspaces"],
        queryFn: async () => {
            const response = await apiClient.get<components["schemas"]["WorkspaceList"]>("/workspaces");
            return response.data.data || [];
        },
    });
}

export function useCreateWorkspace() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newWorkspace: components["schemas"]["WorkspaceCreate"]) => {
            const response = await apiClient.post("/workspaces", newWorkspace);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["workspaces"] });
        },
    });
}
