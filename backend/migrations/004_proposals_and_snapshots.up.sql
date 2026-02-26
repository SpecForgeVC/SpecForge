CREATE TYPE proposal_type AS ENUM ('EDIT_DESCRIPTION', 'MODIFY_SCHEMA', 'ADD_VARIABLE', 'REMOVE_FIELD');
CREATE TYPE proposal_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

CREATE TABLE IF NOT EXISTS ai_proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_item_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    proposal_type proposal_type NOT NULL,
    diff JSONB NOT NULL,
    reasoning TEXT,
    confidence_score FLOAT DEFAULT 0.0,
    status proposal_status DEFAULT 'PENDING',
    reviewed_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS version_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_item_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    snapshot_data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID
);
