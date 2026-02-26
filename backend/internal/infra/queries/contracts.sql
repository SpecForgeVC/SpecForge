-- name: GetContractDefinition :one
SELECT * FROM contract_definitions
WHERE id = $1 LIMIT 1;

-- name: ListContractDefinitions :many
SELECT * FROM contract_definitions
WHERE roadmap_item_id = $1
ORDER BY created_at DESC;

-- name: CreateContractDefinition :one
INSERT INTO contract_definitions (
  roadmap_item_id, contract_type, version, input_schema, output_schema, error_schema, backward_compatible, deprecated_fields
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: DeleteContractDefinition :exec
DELETE FROM contract_definitions
WHERE id = $1;

-- name: ListContractDefinitionsByProject :many
SELECT cd.* FROM contract_definitions cd
JOIN roadmap_items ri ON cd.roadmap_item_id = ri.id
WHERE ri.project_id = $1
ORDER BY cd.created_at DESC;

-- name: UpdateContractDefinition :one
UPDATE contract_definitions SET
  contract_type = $2,
  version = $3,
  input_schema = $4,
  output_schema = $5,
  error_schema = $6,
  backward_compatible = $7,
  deprecated_fields = $8
WHERE id = $1
RETURNING *;
