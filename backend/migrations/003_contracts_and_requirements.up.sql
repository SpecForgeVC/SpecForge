CREATE TABLE IF NOT EXISTS requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_item_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    testable BOOLEAN DEFAULT TRUE,
    acceptance_criteria TEXT,
    order_index INTEGER DEFAULT 0
);

CREATE TYPE contract_type AS ENUM ('REST', 'GRAPHQL', 'CLI', 'INTERNAL_FUNCTION', 'EVENT');

CREATE TABLE IF NOT EXISTS contract_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_item_id UUID NOT NULL REFERENCES roadmap_items(id) ON DELETE CASCADE,
    contract_type contract_type NOT NULL,
    version TEXT NOT NULL,
    input_schema JSONB DEFAULT '{}',
    output_schema JSONB DEFAULT '{}',
    error_schema JSONB DEFAULT '{}',
    backward_compatible BOOLEAN DEFAULT TRUE,
    deprecated_fields JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS variable_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL REFERENCES contract_definitions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    required BOOLEAN DEFAULT FALSE,
    default_value TEXT,
    description TEXT,
    validation_rules JSONB DEFAULT '{}'
);
