-- name: GetWorkspace :one
SELECT * FROM workspaces
WHERE id = $1 LIMIT 1;

-- name: ListWorkspaces :many
SELECT * FROM workspaces
ORDER BY name;

-- name: CreateWorkspace :one
INSERT INTO workspaces (
  name, description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateWorkspace :one
UPDATE workspaces
  set name = $2,
  description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces
WHERE id = $1;
