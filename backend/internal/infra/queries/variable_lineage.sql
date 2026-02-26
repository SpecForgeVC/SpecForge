-- name: CreateVariableLineageEvent :one
INSERT INTO variable_lineage_events (
    variable_id,
    event_type,
    source_component,
    description,
    performed_by,
    metadata
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetVariableLineageEvents :many
SELECT * FROM variable_lineage_events
WHERE variable_id = $1
ORDER BY created_at DESC;

-- name: CreateVariableDependency :one
INSERT INTO variable_dependencies (
    source_variable_id,
    target_variable_id,
    dependency_type
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetVariableDependencies :many
SELECT * FROM variable_dependencies
WHERE source_variable_id = $1 OR target_variable_id = $1;

-- name: DeleteVariableDependencies :exec
DELETE FROM variable_dependencies
WHERE source_variable_id = $1 OR target_variable_id = $1;
