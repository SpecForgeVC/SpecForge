-- Up Migration

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Feature Intelligence
CREATE TABLE IF NOT EXISTS feature_intelligence (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    feature_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    completeness_score INTEGER NOT NULL,
    contract_integrity_score INTEGER NOT NULL,
    variable_coverage_score INTEGER NOT NULL,
    dependency_stability_score INTEGER NOT NULL,
    drift_risk_score INTEGER NOT NULL,
    test_coverage_score INTEGER NOT NULL,
    llm_confidence_score INTEGER NOT NULL,
    overall_score INTEGER NOT NULL,
    last_calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_feature_intelligence UNIQUE (feature_id)
);

CREATE INDEX IF NOT EXISTS idx_feature_intelligence_feature_id ON feature_intelligence(feature_id);

-- 2. Variable Lineage Events
DO $$ BEGIN
    CREATE TYPE lineage_event_type AS ENUM (
        'DECLARED', 'MUTATED', 'TYPE_CHANGED', 'MAPPED_TO_CONTRACT', 'PASSED_TO_API', 'USED_IN_TEST', 'REMOVED'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS variable_lineage_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    variable_id UUID NOT NULL REFERENCES variable_definitions(id) ON DELETE CASCADE,
    event_type lineage_event_type NOT NULL,
    source_component VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    performed_by UUID NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_variable_lineage_variable_id ON variable_lineage_events(variable_id);
CREATE INDEX IF NOT EXISTS idx_variable_lineage_created_at ON variable_lineage_events(created_at);

-- 3. Variable Dependencies
DO $$ BEGIN
    CREATE TYPE dependency_type AS ENUM ('DIRECT', 'DERIVED', 'CONTRACT');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS variable_dependencies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_variable_id UUID NOT NULL REFERENCES variable_definitions(id) ON DELETE CASCADE,
    target_variable_id UUID NOT NULL REFERENCES variable_definitions(id) ON DELETE CASCADE,
    dependency_type dependency_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_variable_dependency UNIQUE (source_variable_id, target_variable_id)
);

CREATE INDEX IF NOT EXISTS idx_variable_dependencies_source ON variable_dependencies(source_variable_id);
CREATE INDEX IF NOT EXISTS idx_variable_dependencies_target ON variable_dependencies(target_variable_id);

-- 4. Add Readiness Column to Roadmap Items
ALTER TABLE roadmap_items ADD COLUMN IF NOT EXISTS readiness_level VARCHAR(50) DEFAULT 'NEEDS_REFINEMENT';
