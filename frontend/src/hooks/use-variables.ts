import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useVariables(projectId?: string) {
    return useQuery({
        queryKey: ["variables", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["VariableList"]>(`/projects/${projectId}/variables`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}

export function useCreateVariable(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newVariable: components["schemas"]["VariableCreateByProject"]) => {
            const response = await apiClient.post(`/projects/${projectId}/variables`, newVariable);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["variables", projectId] });
        },
    });
}

export function useUpdateVariable(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, updates }: { id: string; updates: components["schemas"]["VariableUpdate"] }) => {
            const response = await apiClient.patch(`/variables/${id}`, updates);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["variables", projectId] });
        },
    });
}

export function useDeleteVariable(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            await apiClient.delete(`/variables/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["variables", projectId] });
        },
    });
}
