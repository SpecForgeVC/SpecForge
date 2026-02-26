-- 1. Roadmap Dependencies
CREATE TABLE IF NOT EXISTS roadmap_dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    target_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    dependency_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_roadmap_dependency UNIQUE (source_id, target_id),
    CONSTRAINT no_self_dependency CHECK (source_id != target_id)
);

CREATE INDEX IF NOT EXISTS idx_roadmap_dependencies_source ON roadmap_dependencies(source_id);
CREATE INDEX IF NOT EXISTS idx_roadmap_dependencies_target ON roadmap_dependencies(target_id);

-- 2. Alignment Reports
CREATE TABLE IF NOT EXISTS alignment_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    conflicts JSONB DEFAULT '[]',
    "overlaps" JSONB DEFAULT '[]',
    missing_dependencies JSONB DEFAULT '[]',
    circular_dependencies JSONB DEFAULT '[]',
    recommended_resolutions JSONB DEFAULT '[]',
    alignment_score INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_alignment_reports_project_id ON alignment_reports(project_id);
CREATE INDEX IF NOT EXISTS idx_alignment_reports_created_at ON alignment_reports(created_at);

-- 3. Add alignment_score to projects if not exists
ALTER TABLE projects ADD COLUMN IF NOT EXISTS alignment_score INTEGER DEFAULT 100;
