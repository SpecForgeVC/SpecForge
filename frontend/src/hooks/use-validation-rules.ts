import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useValidationRules(projectId?: string) {
    return useQuery({
        queryKey: ["validation-rules", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["ValidationRuleList"]>(`/projects/${projectId}/validation-rules`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}

export function useCreateValidationRule(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newRule: components["schemas"]["ValidationRuleCreate"]) => {
            const response = await apiClient.post(`/projects/${projectId}/validation-rules`, newRule);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["validation-rules", projectId] });
        },
    });
}

export function useUpdateValidationRule(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, updates }: { id: string; updates: components["schemas"]["ValidationRuleUpdate"] }) => {
            const response = await apiClient.patch(`/validation-rules/${id}`, updates);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["validation-rules", projectId] });
        },
    });
}

export function useDeleteValidationRule(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            await apiClient.delete(`/validation-rules/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["validation-rules", projectId] });
        },
    });
}
