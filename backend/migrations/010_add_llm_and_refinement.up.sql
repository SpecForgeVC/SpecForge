DO $$ BEGIN
    CREATE TYPE llm_provider AS ENUM ('openai', 'ollama', 'gemini', 'anthropic');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE refinement_status AS ENUM ('IN_PROGRESS', 'VALIDATED', 'FAILED', 'APPROVED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS llm_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider llm_provider NOT NULL,
    api_key TEXT NOT NULL,
    base_url TEXT,
    model TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Ensure only one active config
DROP INDEX IF EXISTS one_active_config;
CREATE UNIQUE INDEX one_active_config ON llm_configurations (is_active) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS refinement_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_type TEXT NOT NULL,
    initial_prompt TEXT NOT NULL,
    context_data JSONB,
    max_iterations INTEGER NOT NULL,
    current_iteration INTEGER NOT NULL DEFAULT 0,
    status refinement_status NOT NULL DEFAULT 'IN_PROGRESS',
    confidence_score DOUBLE PRECISION,
    validation_errors TEXT[],
    result JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refinement_iterations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES refinement_sessions(id) ON DELETE CASCADE,
    iteration INTEGER NOT NULL,
    prompt TEXT NOT NULL,
    response TEXT NOT NULL,
    artifact JSONB,
    validation_result JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
