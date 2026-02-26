import { apiClient } from './client';

export interface UIRoadmapItem {
    id: string;
    project_id: string;
    linked_feature_id?: string;
    name: string;
    description: string;
    user_persona: string;
    use_case: string;
    screen_type: string;
    layout_definition: any;
    component_tree: any;
    state_machine: any;
    backend_bindings: any;
    accessibility_spec: any;
    responsive_spec: any;
    validation_rules: any;
    animation_rules: any;
    design_tokens_used: string[];
    edge_cases: any;
    test_scenarios: any;
    intelligence_score: number;
    version: number;
    created_at: string;
    updated_at: string;
}

export interface ExportBundle {
    json_spec: string;
    llm_prompt: string;
    figma_make: string;
    claude_figma: string;
    storybook_spec: string;
}

export const uiRoadmapApi = {
    list: async (projectId: string): Promise<UIRoadmapItem[]> => {
        const response = await apiClient.get(`/projects/${projectId}/ui-roadmap`);
        return response.data;
    },
    get: async (id: string): Promise<UIRoadmapItem> => {
        const response = await apiClient.get(`/ui-roadmap/${id}`);
        return response.data;
    },
    save: async (projectId: string, item: Partial<UIRoadmapItem>): Promise<UIRoadmapItem> => {
        const response = await apiClient.post(`/projects/${projectId}/ui-roadmap`, item);
        return response.data;
    },
    update: async (id: string, item: Partial<UIRoadmapItem>): Promise<UIRoadmapItem> => {
        const response = await apiClient.put(`/ui-roadmap/${id}`, item);
        return response.data;
    },
    delete: async (id: string): Promise<void> => {
        await apiClient.delete(`/ui-roadmap/${id}`);
    },
    export: async (id: string): Promise<ExportBundle> => {
        const response = await apiClient.get(`/ui-roadmap/${id}/export`);
        return response.data;
    },
    sync: async (id: string, payload: any): Promise<void> => {
        await apiClient.post(`/ui-roadmap/${id}/sync`, payload);
    },
    getPluginAssets: async (id: string): Promise<Record<string, string>> => {
        const response = await apiClient.get(`/ui-roadmap/${id}/plugin-assets`);
        return response.data;
    },
};
