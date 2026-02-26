-- name: InsertIntelligenceSnapshot :one
INSERT INTO project_intelligence_snapshots (
  project_id, version, snapshot_json, architecture_score, contract_density, risk_score, alignment_score, confidence_json
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListSnapshotsByProject :many
SELECT * FROM project_intelligence_snapshots
WHERE project_id = $1
ORDER BY version DESC;

-- name: GetIntelligenceSnapshot :one
SELECT * FROM project_intelligence_snapshots
WHERE id = $1 LIMIT 1;

-- name: GetLatestSnapshot :one
SELECT * FROM project_intelligence_snapshots
WHERE project_id = $1
ORDER BY version DESC
LIMIT 1;

-- name: GetMaxSnapshotVersion :one
SELECT COALESCE(MAX(version), 0)::int AS max_version
FROM project_intelligence_snapshots
WHERE project_id = $1;

-- name: InsertProjectModule :one
INSERT INTO project_modules (
  project_id, snapshot_id, name, description, risk_level, change_sensitivity
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListModulesBySnapshot :many
SELECT * FROM project_modules
WHERE snapshot_id = $1
ORDER BY name;

-- name: InsertProjectEntity :one
INSERT INTO project_entities (
  project_id, snapshot_id, name, relationships_json, constraints_json
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListEntitiesBySnapshot :many
SELECT * FROM project_entities
WHERE snapshot_id = $1
ORDER BY name;

-- name: InsertProjectApiEntry :one
INSERT INTO project_api_index (
  project_id, snapshot_id, endpoint, method, auth_type, request_schema, response_schema
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListApiEntriesBySnapshot :many
SELECT * FROM project_api_index
WHERE snapshot_id = $1
ORDER BY endpoint;

-- name: InsertProjectContractEntry :one
INSERT INTO project_contract_registry (
  project_id, snapshot_id, name, contract_type, schema_json, source_module, stability_score
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListContractEntriesBySnapshot :many
SELECT * FROM project_contract_registry
WHERE snapshot_id = $1
ORDER BY name;
