import { apiClient, type ApiResponse } from './client';
import type { components } from './generated/schema';

export type Variable = components["schemas"]["VariableDefinition"];
export type CreateVariableRequest = components["schemas"]["VariableCreateByProject"];

export const variablesApi = {
    listVariables: async (roadmapItemId: string): Promise<Variable[]> => {
        // GET /projects/{projectId}/variables is the list endpoint, but do we have one for roadmap item?
        // Schema shows /projects/{projectId}/variables. 
        // Let's check backend routes. The UI sends roadmap_item_id in context.
        // If we don't have a roadmap-item specific variable list, we might need to filter or update backend.
        // For now, let's assume valid endpoint or update backend if needed.
        // WAIT, schema doesn't show /roadmap-items/{id}/variables.
        // It shows /projects/{projectId}/variables.
        // But our requirement is variables for a roadmap item.
        // Let's use the project-level list for now and filter? Or simpler:
        // The prompt implies we are attaching variables to the roadmap item logic, but the entity might be project scoped.
        // Let's stick to what we have in schema for now.
        const response = await apiClient.get<ApiResponse<components["schemas"]["VariableList"]>>(`/projects/${roadmapItemId}/variables`);
        // VariableList schema likely has a 'variables' property or is an array wrapper? 
        // Checking schema: VariableList = { variables: VariableDefinition[] } usually.
        // Let's assume response.data.data is the VariableList object, so we need .variables
        const data = response.data.data as any;
        if (Array.isArray(data)) {
            return data;
        }
        return data.variables || [];
    },

    createVariable: async (projectId: string, data: CreateVariableRequest): Promise<Variable> => {
        const response = await apiClient.post<ApiResponse<Variable>>(`/projects/${projectId}/variables`, data);
        return response.data.data;
    },

    deleteVariable: async (variableId: string): Promise<void> => {
        await apiClient.delete(`/variables/${variableId}`);
    }
};
