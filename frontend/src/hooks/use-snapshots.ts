import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useSnapshots(projectId?: string) {
    return useQuery({
        queryKey: ["snapshots", projectId],
        queryFn: async () => {
            if (!projectId) return [];
            const response = await apiClient.get<components["schemas"]["SnapshotList"]>(`/projects/${projectId}/snapshots`);
            return response.data.data || [];
        },
        enabled: !!projectId,
    });
}
