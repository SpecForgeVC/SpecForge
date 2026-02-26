import { apiClient } from './client';
import type { AlignmentReport, RoadmapDependency } from '../types/alignment';

export const alignmentApi = {
    getAlignmentReport: async (projectId: string): Promise<AlignmentReport> => {
        const response = await apiClient.get(`/projects/${projectId}/alignment`);
        return response.data;
    },

    triggerAlignmentCheck: async (projectId: string): Promise<AlignmentReport> => {
        const response = await apiClient.post(`/projects/${projectId}/alignment`);
        return response.data;
    },

    listRoadmapDependencies: async (projectId: string): Promise<{ success: boolean, data: RoadmapDependency[] }> => {
        const response = await apiClient.get(`/projects/${projectId}/roadmap-dependencies`);
        return response.data;
    },

    createRoadmapDependency: async (projectId: string, dependency: Partial<RoadmapDependency>): Promise<RoadmapDependency> => {
        const response = await apiClient.post(`/projects/${projectId}/roadmap-dependencies`, dependency);
        return response.data;
    },

    deleteRoadmapDependency: async (dependencyId: string): Promise<void> => {
        await apiClient.delete(`/roadmap-dependencies/${dependencyId}`);
    }
};
