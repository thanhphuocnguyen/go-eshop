// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: sessions.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getSession = `-- name: GetSession :one
SELECT id, user_id, refresh_token, user_agent, client_ip, blocked, expired_at, created_at FROM sessions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.Blocked,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSessionByRefreshToken = `-- name: GetSessionByRefreshToken :one
SELECT id, user_id, refresh_token, user_agent, client_ip, blocked, expired_at, created_at FROM sessions
WHERE refresh_token = $1 LIMIT 1
`

func (q *Queries) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	row := q.db.QueryRow(ctx, getSessionByRefreshToken, refreshToken)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.Blocked,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}

const insertSession = `-- name: InsertSession :one
INSERT INTO sessions (
    id,
    user_id,
    refresh_token,
    user_agent,
    client_ip,
    blocked,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, user_id, refresh_token, user_agent, client_ip, blocked, expired_at, created_at
`

type InsertSessionParams struct {
	ID           uuid.UUID          `json:"id"`
	UserID       uuid.UUID          `json:"userId"`
	RefreshToken string             `json:"refreshToken"`
	UserAgent    string             `json:"userAgent"`
	ClientIp     string             `json:"clientIp"`
	Blocked      bool               `json:"blocked"`
	ExpiredAt    pgtype.Timestamptz `json:"expiredAt"`
}

func (q *Queries) InsertSession(ctx context.Context, arg InsertSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, insertSession,
		arg.ID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.Blocked,
		arg.ExpiredAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.Blocked,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateSession = `-- name: UpdateSession :one
UPDATE sessions
SET
    user_agent = COALESCE($2, user_agent),
    client_ip = COALESCE($3, client_ip),
    blocked = COALESCE($4, blocked),
    expired_at = COALESCE($5, expired_at)
WHERE id = $1
RETURNING id, user_id, refresh_token, user_agent, client_ip, blocked, expired_at, created_at
`

type UpdateSessionParams struct {
	ID        uuid.UUID          `json:"id"`
	UserAgent *string            `json:"userAgent"`
	ClientIp  *string            `json:"clientIp"`
	Blocked   *bool              `json:"blocked"`
	ExpiredAt pgtype.Timestamptz `json:"expiredAt"`
}

func (q *Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, updateSession,
		arg.ID,
		arg.UserAgent,
		arg.ClientIp,
		arg.Blocked,
		arg.ExpiredAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.Blocked,
		&i.ExpiredAt,
		&i.CreatedAt,
	)
	return i, err
}
