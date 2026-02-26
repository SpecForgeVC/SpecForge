export interface Project {
    id: string;
    description?: string;
    workspace_id: string;
    name: string;
    tech_stack?: Record<string, unknown>;
    created_at?: string;
    updated_at?: string;
    repository_url?: string;
}

export interface ProjectCreate {
    name: string;
    description: string;
    tech_stack: Record<string, unknown>;
    repository_url: string;
    project_type: "NEW" | "EXISTING";
}

export interface ProjectList {
    data: Project[];
}
