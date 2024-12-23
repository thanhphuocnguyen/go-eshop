// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
    users (
        email,
        username,
        phone,
        full_name,
        hashed_password,
        role
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING id, email, username, full_name, role, verified_email, verified_phone, created_at, updated_at
`

type CreateUserParams struct {
	Email          string   `json:"email"`
	Username       string   `json:"username"`
	Phone          string   `json:"phone"`
	FullName       string   `json:"full_name"`
	HashedPassword string   `json:"hashed_password"`
	Role           UserRole `json:"role"`
}

type CreateUserRow struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Role          UserRole  `json:"role"`
	VerifiedEmail bool      `json:"verified_email"`
	VerifiedPhone bool      `json:"verified_phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.Username,
		arg.Phone,
		arg.FullName,
		arg.HashedPassword,
		arg.Role,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.FullName,
		&i.Role,
		&i.VerifiedEmail,
		&i.VerifiedPhone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    id, role, username, email, phone, full_name, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
FROM
    users
WHERE
    email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.FullName,
		&i.HashedPassword,
		&i.VerifiedEmail,
		&i.VerifiedPhone,
		&i.PasswordChangedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT
    id, role, username, email, phone, full_name, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
FROM
    users
WHERE
    id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.FullName,
		&i.HashedPassword,
		&i.VerifiedEmail,
		&i.VerifiedPhone,
		&i.PasswordChangedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT
    id, role, username, email, phone, full_name, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
FROM
    users
WHERE
    username = $1
LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.FullName,
		&i.HashedPassword,
		&i.VerifiedEmail,
		&i.VerifiedPhone,
		&i.PasswordChangedAt,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT
    id, role, username, email, phone, full_name, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
FROM
    users
ORDER BY
    id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Role,
			&i.Username,
			&i.Email,
			&i.Phone,
			&i.FullName,
			&i.HashedPassword,
			&i.VerifiedEmail,
			&i.VerifiedPhone,
			&i.PasswordChangedAt,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE
    users
SET
    email = coalesce($1, email),
    full_name = coalesce($2, full_name),
    role = coalesce($3, role),
    verified_email = coalesce($4, verified_email),
    verified_phone = coalesce($5, verified_phone),
    hashed_password = coalesce($6, hashed_password),
    password_changed_at = coalesce($7, password_changed_at),
    updated_at = $8
WHERE
    id = $9
RETURNING id, email, username, full_name, role, verified_email, verified_phone, created_at, updated_at
`

type UpdateUserParams struct {
	Email             pgtype.Text        `json:"email"`
	FullName          pgtype.Text        `json:"full_name"`
	Role              NullUserRole       `json:"role"`
	VerifiedEmail     pgtype.Bool        `json:"verified_email"`
	VerifiedPhone     pgtype.Bool        `json:"verified_phone"`
	HashedPassword    pgtype.Text        `json:"hashed_password"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	ID                int64              `json:"id"`
}

type UpdateUserRow struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Role          UserRole  `json:"role"`
	VerifiedEmail bool      `json:"verified_email"`
	VerifiedPhone bool      `json:"verified_phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Email,
		arg.FullName,
		arg.Role,
		arg.VerifiedEmail,
		arg.VerifiedPhone,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.UpdatedAt,
		arg.ID,
	)
	var i UpdateUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.FullName,
		&i.Role,
		&i.VerifiedEmail,
		&i.VerifiedPhone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
