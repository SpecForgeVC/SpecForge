-- name: CreateWebhook :one
INSERT INTO webhooks (
    project_id, url, events, secret, active
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetWebhook :one
SELECT * FROM webhooks WHERE id = $1;

-- name: ListWebhooksByProject :many
SELECT * FROM webhooks WHERE project_id = $1;

-- name: UpdateWebhook :one
UPDATE webhooks SET
    url = $2,
    events = $3,
    secret = $4,
    active = $5
WHERE id = $1
RETURNING *;

-- name: DeleteWebhook :exec
DELETE FROM webhooks WHERE id = $1;
