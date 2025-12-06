-- name: CreateLink :one
INSERT INTO links (profile_id, title, url, position)
VALUES ($1, $2, $3, COALESCE($4, 0))
RETURNING *;

-- name: GetLinksByProfileID :many
SELECT * FROM links 
WHERE profile_id = $1 
ORDER BY position ASC;

-- name: GetLinkByID :one
SELECT * FROM links WHERE id = $1 LIMIT 1;

-- name: UpdateLink :one
UPDATE links
SET title = COALESCE($2, title),
    url = COALESCE($3, url),
    is_active = COALESCE($4, is_active),
    scheduled_at = COALESCE($5, scheduled_at),
    expires_at = COALESCE($6, expires_at),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteLink :exec
DELETE FROM links WHERE id = $1;

-- name: UpdateLinkPosition :exec
UPDATE links SET position = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: IncrementLinkClicks :exec
UPDATE links SET clicks = clicks + 1 WHERE id = $1;
