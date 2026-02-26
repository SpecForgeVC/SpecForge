export type Severity = "INFO" | "WARNING" | "ERROR" | "CRITICAL";

export type ConflictType =
    | "SCHEMA_MISMATCH"
    | "CONTRACT_COLLISION"
    | "LOGIC_CONTRADICTION"
    | "DEPENDENCY_LOOP";

export interface Conflict {
    id: string;
    severity: Severity;
    type: ConflictType;
    source_id: string;
    target_id: string;
    description: string;
    remediation: string;
    created_at: string;
}

export interface Overlap {
    type: string;
    shared_fields: string[];
    description: string;
}

export interface AlignmentReport {
    id: string;
    project_id: string;
    conflicts: Conflict[];
    overlaps: Overlap[];
    missing_dependencies: string[];
    circular_dependencies: string[];
    recommended_resolutions: string[];
    alignment_score: number;
    created_at: string;
}

export type RoadmapDependencyType = "DIRECT" | "DERIVED" | "CONTRACT";

export interface RoadmapDependency {
    id: string;
    source_id: string;
    target_id: string;
    dependency_type: RoadmapDependencyType;
    created_at: string;
}

export interface RoadmapDependencyCreate {
    source_id: string;
    target_id: string;
    dependency_type: RoadmapDependencyType;
}
