import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useRoadmapItem(roadmapItemId?: string) {
    return useQuery({
        queryKey: ["roadmap-item", roadmapItemId],
        queryFn: async () => {
            if (!roadmapItemId) return null;
            const response = await apiClient.get<components["schemas"]["RoadmapItem"]>(`/roadmap-items/${roadmapItemId}`);
            return response.data;
        },
        enabled: !!roadmapItemId,
    });

}

export function useUpdateRoadmapItem() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async ({ id, data }: { id: string, data: Partial<components["schemas"]["RoadmapItem"]> }) => {
            const response = await apiClient.patch<components["schemas"]["RoadmapItem"]>(`/roadmap-items/${id}`, data);
            return response.data;
        },
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ["roadmap-item", data.id] });
        }
    });
}
