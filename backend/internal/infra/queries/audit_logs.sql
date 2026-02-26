-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    entity_type, entity_id, action, performed_by, old_data, new_data
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListAuditLogsByEntity :many
SELECT * FROM audit_logs
WHERE entity_type = $1 AND entity_id = $2
ORDER BY created_at DESC;

-- name: ListAuditLogsByUser :many
SELECT * FROM audit_logs
WHERE performed_by = $1
ORDER BY created_at DESC;

-- name: ListAuditLogsByAction :many
SELECT * FROM audit_logs
WHERE action = $1
ORDER BY created_at DESC;
