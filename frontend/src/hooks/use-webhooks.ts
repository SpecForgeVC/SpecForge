import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useWebhooks(projectId?: string) {
    return useQuery({
        queryKey: ["webhooks", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["WebhookList"]>(`/projects/${projectId}/webhooks`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}

export function useCreateWebhook(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newWebhook: components["schemas"]["WebhookCreate"]) => {
            const response = await apiClient.post(`/projects/${projectId}/webhooks`, newWebhook);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["webhooks", projectId] });
        },
    });
}

export function useUpdateWebhook(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, updates }: { id: string; updates: components["schemas"]["WebhookUpdate"] }) => {
            const response = await apiClient.patch(`/webhooks/${id}`, updates);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["webhooks", projectId] });
        },
    });
}

export function useDeleteWebhook(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            await apiClient.delete(`/webhooks/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["webhooks", projectId] });
        },
    });
}
