-- name: CreateValidationRule :one
INSERT INTO validation_rules (
    project_id, name, rule_type, rule_config, description
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetValidationRule :one
SELECT * FROM validation_rules WHERE id = $1;

-- name: ListValidationRulesByProject :many
SELECT * FROM validation_rules WHERE project_id = $1;

-- name: UpdateValidationRule :one
UPDATE validation_rules SET
    name = $2,
    rule_type = $3,
    rule_config = $4,
    description = $5
WHERE id = $1
RETURNING *;

-- name: DeleteValidationRule :exec
DELETE FROM validation_rules WHERE id = $1;
