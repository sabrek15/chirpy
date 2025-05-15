-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: DeteleUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserCredentials :one
UPDATE users
SET
    email = COALESCE($2, email),
    hashed_password = COALESCE($3, hashed_password),
    updated_at = NOW()
WHERE
    id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpdateUserByID :exec
UPDATE users
SET
    is_chirpy_red = TRUE,
    updated_at = NOW()
WHERE
    id = $1;