-- Add UI Roadmap Items table
CREATE TABLE IF NOT EXISTS ui_roadmap_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    linked_feature_id UUID REFERENCES roadmap_items(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    description TEXT,
    user_persona TEXT,
    use_case TEXT,
    screen_type TEXT, -- page | modal | component | layout
    layout_definition JSONB DEFAULT '{}',
    component_tree JSONB DEFAULT '{}',
    state_machine JSONB DEFAULT '{}',
    backend_bindings JSONB DEFAULT '{}',
    accessibility_spec JSONB DEFAULT '{}',
    responsive_spec JSONB DEFAULT '{}',
    validation_rules JSONB DEFAULT '{}',
    animation_rules JSONB DEFAULT '{}',
    design_tokens_used TEXT[] DEFAULT '{}',
    edge_cases JSONB DEFAULT '{}',
    test_scenarios JSONB DEFAULT '{}',
    intelligence_score DOUBLE PRECISION DEFAULT 0.0,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ui_roadmap_project_id ON ui_roadmap_items(project_id);
CREATE INDEX IF NOT EXISTS idx_ui_roadmap_linked_feature_id ON ui_roadmap_items(linked_feature_id);
