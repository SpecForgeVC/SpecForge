import { apiClient } from './client';

export interface ExportOptions {
    format: 'json' | 'markdown' | 'zip';
    include_dependencies: boolean;
    include_governance: boolean;
}

export const roadmapApi = {
    exportBuildArtifact: async (roadmapItemId: string, options: ExportOptions) => {
        const response = await apiClient.get(`/roadmap-items/${roadmapItemId}/export`, {
            params: options,
            responseType: options.format === 'zip' ? 'blob' : 'json'
        });

        if (options.format === 'zip') {
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `build-artifact-${roadmapItemId}.zip`);
            document.body.appendChild(link);
            link.click();
            link.remove();
        } else {
            // For JSON/Markdown, we can just return the data or trigger download
            const content = options.format === 'json' ? JSON.stringify(response.data, null, 2) : response.data;
            const blob = new Blob([content], { type: options.format === 'json' ? 'application/json' : 'text/markdown' });
            const url = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', `build-artifact-${roadmapItemId}.${options.format}`);
            document.body.appendChild(link);
            link.click();
            link.remove();
        }
    },

    listRoadmapItems: async (projectId: string) => {
        const response = await apiClient.get(`/projects/${projectId}/roadmap-items`);
        return response.data.data;
    }
};
