CREATE TYPE import_status AS ENUM ('partial', 'complete');

CREATE TABLE project_import_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    completeness_score INT NOT NULL DEFAULT 0,
    status import_status NOT NULL DEFAULT 'partial',
    iteration_count INT NOT NULL DEFAULT 0,
    locked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE project_import_artifacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES project_import_sessions(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_project_import_sessions_project_id ON project_import_sessions(project_id);
CREATE INDEX idx_project_import_artifacts_session_id ON project_import_artifacts(session_id);
