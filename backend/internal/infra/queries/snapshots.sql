-- name: GetVersionSnapshot :one
SELECT * FROM version_snapshots
WHERE id = $1 LIMIT 1;

-- name: ListVersionSnapshots :many
SELECT * FROM version_snapshots
WHERE roadmap_item_id = $1
ORDER BY created_at DESC;

-- name: CreateVersionSnapshot :one
INSERT INTO version_snapshots (
  roadmap_item_id, snapshot_data, created_by
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteVersionSnapshot :exec
DELETE FROM version_snapshots
WHERE id = $1;

-- name: ListVersionSnapshotsByProject :many
SELECT vs.* FROM version_snapshots vs
JOIN roadmap_items ri ON vs.roadmap_item_id = ri.id
WHERE ri.project_id = $1
ORDER BY vs.created_at DESC;
