-- name: CreateSession :one
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
) RETURNING *;


-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: UpdateSession :one
UPDATE sessions
SET
    user_agent = COALESCE(sqlc.narg('user_agent'), user_agent),
    client_ip = COALESCE(sqlc.narg('client_ip'), client_ip),
    blocked = COALESCE(sqlc.narg('blocked'), blocked),
    expired_at = COALESCE(sqlc.narg('expired_at'), expired_at)
WHERE id = $1
RETURNING *;