-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (user_id, email, verify_code) VALUES ($1, $2, $3) RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails WHERE user_id = $1 AND email = $2;

-- name: GetVerifyEmailByID :one
SELECT * FROM verify_emails WHERE id = $1;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = $1 AND verify_code = $2 AND expired_at > now()
RETURNING *;