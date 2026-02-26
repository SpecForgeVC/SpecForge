import { apiClient } from './client';

export type ProposalType = 'EDIT_DESCRIPTION' | 'MODIFY_SCHEMA' | 'ADD_VARIABLE' | 'REMOVE_FIELD';
export type ProposalStatus = 'PENDING' | 'APPROVED' | 'REJECTED';

export interface AiProposal {
    id: string;
    roadmap_item_id: string;
    proposal_type: ProposalType;
    diff: Record<string, any>; // Likely { original: string, modified: string } for schema
    reasoning: string;
    confidence_score: number;
    status: ProposalStatus;
    reviewed_by?: string;
    created_at: string;
}

export const proposalsApi = {
    getProposal: async (proposalId: string): Promise<AiProposal> => {
        const response = await apiClient.get<AiProposal>(`/ai-proposals/${proposalId}`);
        return response.data;
    },

    approve: async (proposalId: string): Promise<void> => {
        await apiClient.post(`/ai-proposals/${proposalId}/approve`);
    },

    reject: async (proposalId: string): Promise<void> => {
        await apiClient.post(`/ai-proposals/${proposalId}/reject`);
    }
};
