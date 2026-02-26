import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useRequirements(roadmapItemId?: string) {
    return useQuery({
        queryKey: ["requirements", roadmapItemId],
        queryFn: async () => {
            if (!roadmapItemId) return [];
            const response = await apiClient.get<components["schemas"]["RequirementList"]>(`/roadmap-items/${roadmapItemId}/requirements`);
            return response.data.data || [];
        },
        enabled: !!roadmapItemId,
    });
}

export function useCreateRequirement(roadmapItemId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newReq: components["schemas"]["RequirementCreate"]) => {
            const response = await apiClient.post(`/roadmap-items/${roadmapItemId}/requirements`, newReq);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["requirements", roadmapItemId] });
        },
    });
}

export function useUpdateRequirement(roadmapItemId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ id, updates }: { id: string; updates: components["schemas"]["RequirementUpdate"] }) => {
            const response = await apiClient.patch(`/requirements/${id}`, updates);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["requirements", roadmapItemId] });
        },
    });
}

export function useDeleteRequirement(roadmapItemId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (id: string) => {
            await apiClient.delete(`/requirements/${id}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["requirements", roadmapItemId] });
        },
    });
}
