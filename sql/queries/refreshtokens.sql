-- name: CreateToken :one
INSERT INTO refreshtokens(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: GetUserToken :one
SELECT * FROM refreshtokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refreshtokens
SET revoked_at = NOW() -- Or $2 if you want to pass the timestamp from Go
WHERE token = $1;