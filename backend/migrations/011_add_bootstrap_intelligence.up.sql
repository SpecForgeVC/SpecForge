-- Project Bootstrap Intelligence tables
-- Stores structured intelligence snapshots from IDE-generated codebase analysis

CREATE TABLE IF NOT EXISTS project_intelligence_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    snapshot_json JSONB NOT NULL,
    architecture_score NUMERIC(5,2),
    contract_density NUMERIC(5,2),
    risk_score NUMERIC(5,2),
    alignment_score NUMERIC(5,2),
    confidence_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_intelligence_snapshots_project ON project_intelligence_snapshots(project_id);
CREATE INDEX IF NOT EXISTS idx_intelligence_snapshots_project_version ON project_intelligence_snapshots(project_id, version DESC);

CREATE TABLE IF NOT EXISTS project_modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    snapshot_id UUID NOT NULL REFERENCES project_intelligence_snapshots(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    risk_level TEXT,
    change_sensitivity TEXT
);

CREATE INDEX IF NOT EXISTS idx_project_modules_snapshot ON project_modules(snapshot_id);

CREATE TABLE IF NOT EXISTS project_entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    snapshot_id UUID NOT NULL REFERENCES project_intelligence_snapshots(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    relationships_json JSONB,
    constraints_json JSONB
);

CREATE INDEX IF NOT EXISTS idx_project_entities_snapshot ON project_entities(snapshot_id);

CREATE TABLE IF NOT EXISTS project_api_index (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    snapshot_id UUID NOT NULL REFERENCES project_intelligence_snapshots(id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    auth_type TEXT,
    request_schema JSONB,
    response_schema JSONB
);

CREATE INDEX IF NOT EXISTS idx_project_api_index_snapshot ON project_api_index(snapshot_id);

CREATE TABLE IF NOT EXISTS project_contract_registry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    snapshot_id UUID NOT NULL REFERENCES project_intelligence_snapshots(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    contract_type TEXT,
    schema_json JSONB,
    source_module TEXT,
    stability_score NUMERIC(5,2)
);

CREATE INDEX IF NOT EXISTS idx_project_contract_registry_snapshot ON project_contract_registry(snapshot_id);
