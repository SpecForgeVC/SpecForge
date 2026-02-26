-- Migration to add Reality Anchor Engine (RAE) MCP tables

DO $$ BEGIN
    CREATE TYPE snapshot_state AS ENUM ('initiated', 'awaiting_post', 'analyzing', 'completed', 'failed');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS reality_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    roadmap_item_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    state snapshot_state NOT NULL DEFAULT 'initiated',
    snapshot_json JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS snapshot_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snapshot_id UUID NOT NULL REFERENCES reality_snapshots(id) ON DELETE CASCADE,
    scores_json JSONB NOT NULL,
    verdict TEXT NOT NULL, -- 'approved', 'blocked', 'requires_alignment'
    drift_detected BOOLEAN DEFAULT FALSE,
    alignment_conflicts JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add MCP settings to projects table
ALTER TABLE projects ADD COLUMN IF NOT EXISTS mcp_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS mcp_port INTEGER DEFAULT 8098;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS mcp_bind_address TEXT DEFAULT '0.0.0.0';
ALTER TABLE projects ADD COLUMN IF NOT EXISTS mcp_token_required BOOLEAN DEFAULT TRUE;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS mcp_token TEXT;
