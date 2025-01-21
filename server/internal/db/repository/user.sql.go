// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countUsers = `-- name: CountUsers :one
SELECT
    count(*)
FROM
    users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO
    users (user_id,email,username,phone,fullname,hashed_password,role)
VALUES
    ($1,$2,$3,$4,$5,$6,$7)
RETURNING user_id, email, username, fullname, role, verified_email, verified_phone, created_at, updated_at
`

type CreateUserParams struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Phone          string    `json:"phone"`
	Fullname       string    `json:"fullname"`
	HashedPassword string    `json:"hashed_password"`
	Role           UserRole  `json:"role"`
}

type CreateUserRow struct {
	UserID        uuid.UUID `json:"user_id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Fullname      string    `json:"fullname"`
	Role          UserRole  `json:"role"`
	VerifiedEmail bool      `json:"verified_email"`
	VerifiedPhone bool      `json:"verified_phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.UserID,
		arg.Email,
		arg.Username,
		arg.Phone,
		arg.Fullname,
		arg.HashedPassword,
		arg.Role,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.Fullname,
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
    user_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, userID)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    user_id, role, username, email, phone, fullname, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
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
		&i.UserID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.Fullname,
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
    user_id, role, username, email, phone, fullname, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
FROM
    users
WHERE
    user_id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.Fullname,
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
    user_id, role, username, email, phone, fullname, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
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
		&i.UserID,
		&i.Role,
		&i.Username,
		&i.Email,
		&i.Phone,
		&i.Fullname,
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
    user_id, role, username, email, phone, fullname, hashed_password, verified_email, verified_phone, password_changed_at, updated_at, created_at
FROM
    users
ORDER BY
    user_id
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
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Role,
			&i.Username,
			&i.Email,
			&i.Phone,
			&i.Fullname,
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

type SeedUsersParams struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Phone          string    `json:"phone"`
	Fullname       string    `json:"fullname"`
	HashedPassword string    `json:"hashed_password"`
	Role           UserRole  `json:"role"`
}

const updateUser = `-- name: UpdateUser :one
UPDATE
    users
SET
    email = coalesce($1, email),
    fullname = coalesce($2, fullname),
    role = coalesce($3, role),
    verified_email = coalesce($4, verified_email),
    verified_phone = coalesce($5, verified_phone),
    hashed_password = coalesce($6, hashed_password),
    password_changed_at = coalesce($7, password_changed_at),
    updated_at = $8
WHERE
    user_id = $9
RETURNING user_id, email, username, fullname, role, verified_email, verified_phone, created_at, updated_at
`

type UpdateUserParams struct {
	Email             pgtype.Text        `json:"email"`
	Fullname          pgtype.Text        `json:"fullname"`
	Role              NullUserRole       `json:"role"`
	VerifiedEmail     pgtype.Bool        `json:"verified_email"`
	VerifiedPhone     pgtype.Bool        `json:"verified_phone"`
	HashedPassword    pgtype.Text        `json:"hashed_password"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	ID                uuid.UUID          `json:"id"`
}

type UpdateUserRow struct {
	UserID        uuid.UUID `json:"user_id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Fullname      string    `json:"fullname"`
	Role          UserRole  `json:"role"`
	VerifiedEmail bool      `json:"verified_email"`
	VerifiedPhone bool      `json:"verified_phone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Email,
		arg.Fullname,
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
		&i.UserID,
		&i.Email,
		&i.Username,
		&i.Fullname,
		&i.Role,
		&i.VerifiedEmail,
		&i.VerifiedPhone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
