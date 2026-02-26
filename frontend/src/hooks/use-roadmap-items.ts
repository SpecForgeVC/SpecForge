import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useRoadmapItems(projectId?: string) {
    return useQuery({
        queryKey: ["roadmap-items", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["RoadmapItemList"]>(`/projects/${projectId}/roadmap-items`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}

export function useRoadmapItem(itemId?: string) {
    return useQuery({
        queryKey: ["roadmap-item", itemId],
        queryFn: async () => {
            if (!itemId) return null;
            const response = await apiClient.get<components["schemas"]["RoadmapItem"]>(`/roadmap-items/${itemId}`);
            return response.data;
        },
        enabled: !!itemId,
    });
}
export function useCreateRoadmapItem(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newItem: components["schemas"]["RoadmapItemCreate"]) => {
            const response = await apiClient.post(`/projects/${projectId}/roadmap-items`, newItem);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["roadmap-items", projectId] });
        },
    });
}

export function useUpdateRoadmapItem(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, updates }: { id: string; updates: components["schemas"]["RoadmapItemUpdate"] }) => {
            const response = await apiClient.patch(`/roadmap-items/${id}`, updates);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["roadmap-items", projectId] });
        },
    });
}

export function useDeleteRoadmapItem(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            await apiClient.delete(`/roadmap-items/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["roadmap-items", projectId] });
        },
    });
}
