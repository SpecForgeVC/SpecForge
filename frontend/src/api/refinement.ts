import { apiClient as client } from './client';
import type { components } from './generated/schema';

export type RefinementSession = components["schemas"]["RefinementSession"];
export type StartRefinementSessionRequest = components["schemas"]["StartRefinementSession"];

export interface RefinementEvent {
    type: string;
    message: string;
    payload?: any;
    timestamp: string;
}

export const refinementApi = {
    startSession: async (artifactType: string, targetType: string, prompt: string, contextData: any, maxIterations: number): Promise<RefinementSession> => {
        const response = await client.post('/refinement', {
            artifact_type: artifactType,
            target_type: targetType,
            prompt,
            context_data: contextData,
            max_iterations: maxIterations
        });
        return response.data;
    },

    approveSession: async (sessionId: string): Promise<void> => {
        await client.post(`/refinement/${sessionId}/approve`);
    },

    getEventsEndpoint: (sessionId: string) => `/refinement/${sessionId}/events`,

    getEvents: (sessionId: string): EventSource => {
        return new EventSource(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/refinement/${sessionId}/events`);
    }
};
