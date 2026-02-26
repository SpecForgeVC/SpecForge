import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useProject(projectId?: string) {
    return useQuery({
        queryKey: ["project", projectId],
        queryFn: async () => {
            if (!projectId) return null;
            const response = await apiClient.get<components["schemas"]["Project"]>(`/projects/${projectId}`);
            return response.data;
        },
        enabled: !!projectId,
    });
}
