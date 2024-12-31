-- name: CreateSession :one
INSERT INTO sessions (
    session_id,
    user_id,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expired_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: GetSession :one
SELECT * FROM sessions
WHERE session_id = $1 LIMIT 1;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = $1 LIMIT 1;
