-- name: GetRoadmapItem :one
SELECT * FROM roadmap_items
WHERE id = $1 LIMIT 1;

-- name: ListRoadmapItems :many
SELECT * FROM roadmap_items
WHERE project_id = $1
ORDER BY created_at DESC;

-- name: CreateRoadmapItem :one
INSERT INTO roadmap_items (
  project_id, type, title, description, business_context, technical_context, priority, status, risk_level, readiness_level, breaking_change, regression_sensitive
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateRoadmapItem :one
UPDATE roadmap_items
  set title = $2,
  description = $3,
  business_context = $4,
  technical_context = $5,
  status = $6
WHERE id = $1
RETURNING id, project_id, type, title, description, business_context, technical_context, priority, status, risk_level, readiness_level, breaking_change, regression_sensitive, created_at, updated_at;

-- name: DeleteRoadmapItem :exec
DELETE FROM roadmap_items
WHERE id = $1;
