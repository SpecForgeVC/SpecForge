import { apiClient as client } from './client';

export interface LLMConfiguration {
    id?: string;
    provider: 'openai' | 'ollama' | 'gemini' | 'anthropic';
    api_key: string;
    base_url?: string;
    model: string;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
}

export const llmApi = {
    getConfig: async (): Promise<LLMConfiguration | null> => {
        const response = await client.get('/settings/llm');
        return response.data;
    },

    updateConfig: async (config: LLMConfiguration): Promise<LLMConfiguration> => {
        const response = await client.put('/settings/llm', config);
        return response.data;
    },

    testConnection: async (config: LLMConfiguration): Promise<{ message: string }> => {
        const response = await client.post('/settings/llm/test', config);
        return response.data;
    },

    listModels: async (config: LLMConfiguration): Promise<string[]> => {
        const response = await client.post('/settings/llm/models', config);
        return response.data.models;
    },

    // Warmup is handled via SSE directly in component usually, or we can use EventSource
    getWarmupEndpoint: () => '/settings/llm/warmup',
};
