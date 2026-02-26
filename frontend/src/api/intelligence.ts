import { apiClient, type ApiResponse } from './client';

export interface BreakingChange {
    field: string;
    issue: string;
}

export interface DriftReport {
    drift_detected: boolean;
    breaking_changes: BreakingChange[];
    risk_score: number;
}

export interface AuditLog {
    id: string;
    entity_type: string;
    entity_id: string;
    action: string;
    performed_by: string;
    old_data: any;
    new_data: {
        drift_report?: DriftReport;
        version?: string;
        [key: string]: any;
    };
    created_at: string;
}

export interface FeatureIntelligence {
    id: string;
    feature_id: string;
    completeness_score: number;
    contract_integrity_score: number;
    variable_coverage_score: number;
    dependency_stability_score: number;
    drift_risk_score: number;
    test_coverage_score: number;
    llm_confidence_score: number;
    overall_score: number;
    last_calculated_at: string;
}

export interface DriftFix {
    field: string;
    issue: string;
    suggested_change: string;
    explanation: string;
}

export interface LineageData {
    label: string;
    [key: string]: any;
}

export interface LineageNode {
    id: string;
    type: string; // 'variable'
    data: LineageData;
}

export interface LineageEdge {
    id: string;
    source: string;
    target: string;
}

export interface LineageGraph {
    nodes: LineageNode[];
    edges: LineageEdge[];
}

export const intelligenceApi = {
    // Drift history uses SuccessResponse? No, need to check handle. Assuming it does.
    getDriftHistory: async (): Promise<AuditLog[]> => {
        const response = await apiClient.get<ApiResponse<AuditLog[]>>('/drift/history');
        return response.data.data || [];
    },

    getFeatureIntelligence: async (roadmapItemId: string): Promise<FeatureIntelligence> => {
        const response = await apiClient.get<ApiResponse<FeatureIntelligence>>(`/roadmap-items/${roadmapItemId}/intelligence`);
        return response.data.data;
    },

    getLineageGraph: async (variableId: string): Promise<LineageGraph> => {
        const response = await apiClient.get<ApiResponse<LineageGraph>>(`/variables/${variableId}/lineage`);
        return response.data.data;
    },

    runDriftCheck: async (contractId: string, againstVersion: string): Promise<DriftReport> => {
        const response = await apiClient.post<DriftReport>(`/contracts/${contractId}/drift-check`, {
            against_version: againstVersion,
        });
        return response.data;
    },

    generateDriftFixes: async (driftReport: DriftReport, roadmapItemId: string): Promise<DriftFix[]> => {
        const response = await apiClient.post<{ fixes: DriftFix[] }>('/drift/generate-fixes', {
            drift_report: driftReport,
            roadmap_item_id: roadmapItemId,
        });
        return response.data.fixes;
    },

    getRoadmapItemActivity: async (roadmapItemId: string, limit: number = 20): Promise<AuditLog[]> => {
        const response = await apiClient.get<{ success: boolean; data: AuditLog[] }>(`/roadmap-items/${roadmapItemId}/activity`, {
            params: { limit },
        });
        return response.data.data;
    },
};
