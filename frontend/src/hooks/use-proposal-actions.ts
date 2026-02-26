import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/api/client";

export function useApproveProposal(proposalId: string, roadmapItemId: string) {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async () => {
            await apiClient.post(`/ai-proposals/${proposalId}/approve`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["ai-proposals", roadmapItemId] });
            queryClient.invalidateQueries({ queryKey: ["roadmap-item", roadmapItemId] });
        },
    });
}

export function useRejectProposal(proposalId: string, roadmapItemId: string) {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async () => {
            await apiClient.post(`/ai-proposals/${proposalId}/reject`);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ["ai-proposals", roadmapItemId] });
        },
    });
}
