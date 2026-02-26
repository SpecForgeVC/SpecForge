import { apiClient, type ApiResponse } from './client';

export type ContractType = 'REST' | 'GRAPHQL' | 'CLI' | 'INTERNAL_FUNCTION' | 'EVENT';

export interface ContractDefinition {
    id: string;
    roadmap_item_id: string;
    contract_type: ContractType;
    version: string;
    input_schema: Record<string, any>;
    output_schema: Record<string, any>;
    error_schema: Record<string, any>;
    backward_compatible: boolean;
    deprecated_fields: string[];
    created_at: string;
}

export const contractsApi = {
    listContracts: async (roadmapItemId: string): Promise<ContractDefinition[]> => {
        const response = await apiClient.get<ApiResponse<ContractDefinition[]>>(`/roadmap-items/${roadmapItemId}/contracts`);
        return response.data.data;
    },

    getContract: async (contractId: string): Promise<ContractDefinition> => {
        const response = await apiClient.get<ApiResponse<ContractDefinition>>(`/contracts/${contractId}`);
        return response.data.data;
    },

    updateContract: async (contractId: string, data: Partial<ContractDefinition>): Promise<ContractDefinition> => {
        const response = await apiClient.patch<ApiResponse<ContractDefinition>>(`/contracts/${contractId}`, data);
        return response.data.data;
    },

    createContract: async (projectId: string, data: any): Promise<ContractDefinition> => {
        const response = await apiClient.post<ApiResponse<ContractDefinition>>(`/projects/${projectId}/contracts`, data);
        return response.data.data;
    }
};
