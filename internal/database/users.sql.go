// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const deteleUsers = `-- name: DeteleUsers :exec
DELETE FROM users
`

func (q *Queries) DeteleUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deteleUsers)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserByID = `-- name: UpdateUserByID :exec
UPDATE users
SET
    is_chirpy_red = TRUE,
    updated_at = NOW()
WHERE
    id = $1
`

func (q *Queries) UpdateUserByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, updateUserByID, id)
	return err
}

const updateUserCredentials = `-- name: UpdateUserCredentials :one
UPDATE users
SET
    email = COALESCE($2, email),
    hashed_password = COALESCE($3, hashed_password),
    updated_at = NOW()
WHERE
    id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red
`

type UpdateUserCredentialsParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

type UpdateUserCredentialsRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	IsChirpyRed bool
}

func (q *Queries) UpdateUserCredentials(ctx context.Context, arg UpdateUserCredentialsParams) (UpdateUserCredentialsRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserCredentials, arg.ID, arg.Email, arg.HashedPassword)
	var i UpdateUserCredentialsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}
