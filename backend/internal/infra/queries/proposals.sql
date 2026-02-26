-- name: GetAiProposal :one
SELECT * FROM ai_proposals
WHERE id = $1 LIMIT 1;

-- name: ListAiProposalsByProject :many
SELECT * FROM ai_proposals
WHERE roadmap_item_id IN (
    SELECT id FROM roadmap_items WHERE project_id = $1
)
ORDER BY created_at DESC;

-- name: ListAiProposalsByRoadmapItem :many
SELECT * FROM ai_proposals
WHERE roadmap_item_id = $1
ORDER BY created_at DESC;

-- name: CreateAiProposal :one
INSERT INTO ai_proposals (
    roadmap_item_id, proposal_type, diff, reasoning, confidence_score, status
) VALUES (
    $1, $2, $3, $4, $5, 'PENDING'
)
RETURNING *;

-- name: UpdateAiProposalStatus :one
UPDATE ai_proposals
SET status = $2, reviewed_by = $3
WHERE id = $1
RETURNING *;

-- name: DeleteAiProposal :exec
DELETE FROM ai_proposals
WHERE id = $1;
