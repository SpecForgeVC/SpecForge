-- Migration to add MCP API tokens table for IDE connection

CREATE TABLE IF NOT EXISTS mcp_api_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    token_prefix TEXT NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster lookup by hash (though bcrypt/argon2 usually handles the comparison)
-- We store the hash of the token, but we might want to lookup by project_id first if we use a prefix + secret scheme.
-- For simple 64-char random tokens, we hash them and look them up.
CREATE INDEX idx_mcp_api_tokens_user_id ON mcp_api_tokens(user_id);
CREATE INDEX idx_mcp_api_tokens_project_id ON mcp_api_tokens(project_id);
CREATE INDEX idx_mcp_api_tokens_token_hash ON mcp_api_tokens(token_hash);
