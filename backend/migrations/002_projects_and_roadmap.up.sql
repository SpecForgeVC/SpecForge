CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    tech_stack JSONB DEFAULT '{}',
    repository_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE roadmap_item_type AS ENUM ('EPIC', 'FEATURE', 'TASK', 'BUGFIX', 'REFACTOR');
CREATE TYPE roadmap_item_priority AS ENUM ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL');
CREATE TYPE roadmap_item_status AS ENUM ('DRAFT', 'IN_REVIEW', 'APPROVED', 'IN_PROGRESS', 'COMPLETE');
CREATE TYPE risk_level AS ENUM ('LOW', 'MEDIUM', 'HIGH');

CREATE TABLE IF NOT EXISTS roadmap_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    type roadmap_item_type NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    business_context TEXT,
    technical_context TEXT,
    priority roadmap_item_priority DEFAULT 'MEDIUM',
    status roadmap_item_status DEFAULT 'DRAFT',
    risk_level risk_level DEFAULT 'LOW',
    breaking_change BOOLEAN DEFAULT FALSE,
    regression_sensitive BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_projects_updated_at
BEFORE UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roadmap_items_updated_at
BEFORE UPDATE ON roadmap_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
