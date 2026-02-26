-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    full_name TEXT,
    role TEXT NOT NULL, -- OWNER, ADMIN, REVIEWER, ENGINEER, AI_AGENT
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Trigger for updated_at
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Seed initial workspace
INSERT INTO workspaces (id, name, description)
VALUES ('00000000-0000-0000-0000-000000000001', 'Default Workspace', 'Initial system workspace')
ON CONFLICT (id) DO NOTHING;

-- Seed superadmin user
INSERT INTO users (id, email, full_name, role)
VALUES (
    '11111111-1111-1111-1111-111111111111', 
    'admin@specforge.io', 
    'System Superadmin', 
    'OWNER'
)
ON CONFLICT (email) DO NOTHING;
