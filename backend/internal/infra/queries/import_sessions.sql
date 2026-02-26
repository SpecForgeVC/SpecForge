-- name: CreateImportSession :one
INSERT INTO project_import_sessions (
    project_id,
    completeness_score,
    status,
    iteration_count,
    locked
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetImportSession :one
SELECT * FROM project_import_sessions
WHERE id = $1;

-- name: GetLatestImportSessionByProject :one
SELECT * FROM project_import_sessions
WHERE project_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateImportSession :one
UPDATE project_import_sessions
SET 
    completeness_score = $2,
    status = $3,
    iteration_count = $4,
    locked = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CreateImportArtifact :one
INSERT INTO project_import_artifacts (
    session_id,
    payload
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListImportArtifactsBySession :many
SELECT * FROM project_import_artifacts
WHERE session_id = $1
ORDER BY created_at ASC;
