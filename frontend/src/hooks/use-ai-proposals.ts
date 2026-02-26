import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useAIProposals(roadmapItemId?: string) {
    return useQuery({
        queryKey: ["ai-proposals", roadmapItemId],
        queryFn: async () => {
            if (!roadmapItemId) return [];
            const response = await apiClient.get<components["schemas"]["AIProposalList"]>(`/roadmap-items/${roadmapItemId}/ai-proposals`);
            return response.data.data || [];
        },
        enabled: !!roadmapItemId,
    });
}

export function useProjectProposals(projectId?: string) {
    return useQuery({
        queryKey: ["ai-proposals", "project", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["AIProposalList"]>(`/projects/${projectId}/ai-proposals`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}
