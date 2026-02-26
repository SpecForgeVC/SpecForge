import { apiClient, type ApiResponse } from './client';

// --- Types ---

export interface BootstrapScores {
    architecture_score: number;
    contract_density: number;
    risk_score: number;
    alignment_score: number;
}

import type { components } from './generated/schema';

export type ImportSession = components['schemas']['ImportSession'];

export interface BootstrapConfidence {
    project_overview: number;
    tech_stack: number;
    modules: number;
    apis: number;
    data_models: number;
    validation_rules: number;
    contracts: number;
    current_state: number;
    risks: number;
    change_sensitivity: number;
}

export interface BootstrapModule {
    name: string;
    description: string;
    responsibilities?: string[];
    risk_level: 'LOW' | 'MEDIUM' | 'HIGH' | 'UNKNOWN';
    change_sensitivity: 'LOW' | 'MEDIUM' | 'HIGH' | 'UNKNOWN';
}

export interface BootstrapApiEntry {
    endpoint: string;
    method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
    auth_type: string;
    request_schema?: Record<string, any>;
    response_schema?: Record<string, any>;
}

export interface BootstrapEntity {
    name: string;
    relationships?: Record<string, any>[];
    constraints?: Record<string, any>[];
}

export interface BootstrapContractEntry {
    name: string;
    contract_type: string;
    schema?: Record<string, any>;
    source_module: string;
    stability_score: number;
}

export interface ProjectIntelligenceSnapshot {
    id: string;
    project_id: string;
    version: number;
    snapshot_json: any;
    scores: BootstrapScores;
    confidence: BootstrapConfidence;
    modules?: BootstrapModule[];
    entities?: BootstrapEntity[];
    apis?: BootstrapApiEntry[];
    contracts?: BootstrapContractEntry[];
    created_at: string;
    // Flatten score fields for convenience
    architecture_score: number;
    contract_density: number;
    risk_score: number;
    alignment_score: number;
}

export interface BootstrapPromptResponse {
    prompt: string;
    project_name: string;
    project_id: string;
}

export interface BootstrapIngestResponse {
    snapshot: ProjectIntelligenceSnapshot;
    scores: BootstrapScores;
    confidence: BootstrapConfidence;
    warnings: string[];
}

export interface BootstrapDiffResult {
    from_snapshot_id: string | null;
    to_snapshot_id: string;
    from_version: number;
    to_version: number;
    added_modules?: BootstrapModule[];
    removed_modules?: BootstrapModule[];
    added_apis?: BootstrapApiEntry[];
    removed_apis?: BootstrapApiEntry[];
    added_entities?: BootstrapEntity[];
    removed_entities?: BootstrapEntity[];
    score_changes?: Record<string, number>;
}

export interface BootstrapPayload {
    project_overview?: Record<string, any>;
    tech_stack?: Record<string, any>;
    modules?: Record<string, any>[];
    apis?: Record<string, any>[];
    data_models?: Record<string, any>[];
    validation_rules?: Record<string, any>[];
    contracts?: Record<string, any>[];
    current_state?: Record<string, any>;
    risks?: Record<string, any>[];
    change_sensitivity?: Record<string, any>[];
}

// --- API Functions ---

export const bootstrapApi = {
    generatePrompt: async (projectId: string): Promise<BootstrapPromptResponse> => {
        const response = await apiClient.post<ApiResponse<BootstrapPromptResponse>>(
            `projects/${projectId}/bootstrap/generate-prompt`
        );
        return response.data.data;
    },

    ingestBootstrap: async (projectId: string, payload: BootstrapPayload): Promise<BootstrapIngestResponse> => {
        const response = await apiClient.post<ApiResponse<BootstrapIngestResponse>>(
            `projects/${projectId}/bootstrap/ingest`,
            payload
        );
        return response.data.data;
    },

    listSnapshots: async (projectId: string): Promise<ProjectIntelligenceSnapshot[]> => {
        const response = await apiClient.get<ApiResponse<ProjectIntelligenceSnapshot[]>>(
            `projects/${projectId}/bootstrap/snapshots`
        );
        return response.data.data;
    },

    getSnapshot: async (projectId: string, snapshotId: string): Promise<ProjectIntelligenceSnapshot> => {
        const response = await apiClient.get<ApiResponse<ProjectIntelligenceSnapshot>>(
            `projects/${projectId}/bootstrap/snapshots/${snapshotId}`
        );
        return response.data.data;
    },

    getLatestSnapshot: async (projectId: string): Promise<ProjectIntelligenceSnapshot> => {
        const response = await apiClient.get<ApiResponse<ProjectIntelligenceSnapshot>>(
            `projects/${projectId}/bootstrap/latest`
        );
        return response.data.data;
    },

    diffSnapshots: async (
        projectId: string,
        fromSnapshotId?: string,
        toSnapshotId?: string
    ): Promise<BootstrapDiffResult> => {
        const response = await apiClient.post<ApiResponse<BootstrapDiffResult>>(
            `projects/${projectId}/bootstrap/diff`,
            {
                from_snapshot_id: fromSnapshotId || undefined,
                to_snapshot_id: toSnapshotId || undefined,
            }
        );
        return response.data.data;
    },

    getLatestImportSession: async (projectId: string): Promise<ImportSession> => {
        const response = await apiClient.get<ApiResponse<ImportSession>>(
            `projects/${projectId}/bootstrap/session`
        );
        return response.data.data;
    },

    createProject: async (data: { name: string }): Promise<{ id: string; name: string }> => {
        const response = await apiClient.post<ApiResponse<{ id: string; name: string }>>(
            'projects',
            data
        );
        return response.data.data;
    },
};
