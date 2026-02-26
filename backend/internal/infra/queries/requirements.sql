-- name: CreateRequirement :one
INSERT INTO requirements (
    roadmap_item_id, title, description, testable, acceptance_criteria, order_index
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetRequirement :one
SELECT * FROM requirements WHERE id = $1;

-- name: ListRequirementsByRoadmapItem :many
SELECT * FROM requirements WHERE roadmap_item_id = $1 ORDER BY order_index ASC;

-- name: UpdateRequirement :one
UPDATE requirements SET
    title = $2,
    description = $3,
    testable = $4,
    acceptance_criteria = $5,
    order_index = $6
WHERE id = $1
RETURNING *;

-- name: DeleteRequirement :exec
DELETE FROM requirements WHERE id = $1;
