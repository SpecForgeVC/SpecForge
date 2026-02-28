import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient, type ApiResponse } from "@/api/client";
import type { components } from "@/api/generated/schema";

export function useProjects(workspaceId?: string) {
    return useQuery({
        queryKey: ["projects", workspaceId],
        queryFn: async () => {
            if (!workspaceId) return [];
            const response = await apiClient.get<components["schemas"]["ProjectList"]>(`/workspaces/${workspaceId}/projects`);
            return response.data.data || [];
        },
        enabled: !!workspaceId,
    });
}

export function useCreateProject(workspaceId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newProject: components["schemas"]["ProjectCreate"]) => {
            const response = await apiClient.post(`/workspaces/${workspaceId}/projects`, newProject);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["projects", workspaceId] });
        },
    });
}

export function useUpdateProject(projectId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (update: components["schemas"]["ProjectUpdate"]) => {
            const response = await apiClient.patch(`/projects/${projectId}`, update);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["projects"] });
            queryClient.invalidateQueries({ queryKey: ["project", projectId] });
        },
    });
}

export function useTechStackRecommendation() {
    return useMutation({
        mutationFn: async (req: { purpose: string; constraints?: string }) => {
            const response = await apiClient.post<ApiResponse<{ recommended_stack: any; reasoning: string }>>(
                `/projects/recommend-stack`,
                req
            );
            return response.data.data;
        },
    });
}
export function useDeleteProject(workspaceId: string) {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (projectId: string) => {
            await apiClient.delete(`/projects/${projectId}`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["projects", workspaceId] });
        },
    });
}
