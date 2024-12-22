-- name: CreateUser :one
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
RETURNING id, email, username, full_name, role, verified_email, verified_phone, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = coalesce(sqlc.narg('email'), email),
    full_name = coalesce(sqlc.narg('full_name'), full_name),
    role = coalesce(sqlc.narg('role'), role),
    verified_email = coalesce(sqlc.narg('verified_email'), verified_email),
    verified_phone = coalesce(sqlc.narg('verified_phone'), verified_phone),
    hashed_password = coalesce(sqlc.narg('hashed_password'), hashed_password),
    password_changed_at = coalesce(sqlc.narg('password_changed_at'), password_changed_at),
    updated_at = sqlc.arg('updated_at')
WHERE
    id = sqlc.arg('id')
RETURNING id, email, username, full_name, role, verified_email, verified_phone, created_at, updated_at;


-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    id = $1;