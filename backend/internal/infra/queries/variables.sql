-- name: CreateVariable :one
INSERT INTO variable_definitions (
    contract_id, name, type, required, default_value, description, validation_rules
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetVariable :one
SELECT * FROM variable_definitions WHERE id = $1;

-- name: ListVariablesByContract :many
SELECT * FROM variable_definitions WHERE contract_id = $1;

-- name: UpdateVariable :one
UPDATE variable_definitions SET
    name = $2,
    type = $3,
    required = $4,
    default_value = $5,
    description = $6,
    validation_rules = $7
WHERE id = $1
RETURNING *;

-- name: DeleteVariable :exec
DELETE FROM variable_definitions WHERE id = $1;

-- name: ListVariablesByProject :many
SELECT vd.* FROM variable_definitions vd
JOIN contract_definitions cd ON vd.contract_id = cd.id
JOIN roadmap_items ri ON cd.roadmap_item_id = ri.id
WHERE ri.project_id = $1;
