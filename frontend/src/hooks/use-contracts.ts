import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useContracts(projectId?: string) {
    return useQuery({
        queryKey: ["contracts", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["ContractList"]>(`/projects/${projectId}/contracts`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}

export function useCreateContract(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newContract: components["schemas"]["ContractCreateByProject"]) => {
            const response = await apiClient.post(`/projects/${projectId}/contracts`, newContract);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["contracts", projectId] });
        },
    });
}

export function useUpdateContract(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, updates }: { id: string; updates: components["schemas"]["ContractUpdate"] }) => {
            const response = await apiClient.patch(`/contracts/${id}`, updates);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["contracts", projectId] });
        },
    });
}

export function useDeleteContract(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            await apiClient.delete(`/contracts/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["contracts", projectId] });
        },
    });
}
