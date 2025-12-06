-- name: CreateProfile :one
INSERT INTO profiles (user_id)
VALUES ($1)
RETURNING *;

-- name: GetProfileByUserID :one
SELECT * FROM profiles WHERE user_id = $1 LIMIT 1;

-- name: GetProfileByUsername :one
SELECT p.*, u.username 
FROM profiles p
JOIN users u ON p.user_id = u.id
WHERE u.username = $1 LIMIT 1;

-- name: UpdateProfile :one
UPDATE profiles
SET avatar_url = COALESCE($2, avatar_url),
    bio = COALESCE($3, bio),
    theme_config = COALESCE($4, theme_config),
    custom_css = COALESCE($5, custom_css),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING *;
