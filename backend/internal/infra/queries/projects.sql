-- name: GetProject :one
SELECT * FROM projects
WHERE id = $1 LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
WHERE workspace_id = $1
ORDER BY name;

-- name: CreateProject :one
INSERT INTO projects (
  workspace_id, name, description, tech_stack, settings, repository_url,
  mcp_enabled, mcp_port, mcp_bind_address, mcp_token_required, mcp_token
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
  set name = $2,
  description = $3,
  tech_stack = $4,
  settings = $5,
  repository_url = $6,
  mcp_enabled = $7,
  mcp_port = $8,
  mcp_bind_address = $9,
  mcp_token_required = $10,
  mcp_token = $11
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;
