-- name: CreateUser :one
INSERT INTO
    users (
        email,
        username,
        phone,
        fullname,
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
RETURNING user_id, email, username, fullname, role, verified_email, verified_phone, created_at, updated_at;

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
    user_id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT
    *
FROM
    users
ORDER BY
    user_id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = coalesce(sqlc.narg('email'), email),
    fullname = coalesce(sqlc.narg('fullname'), fullname),
    role = coalesce(sqlc.narg('role'), role),
    verified_email = coalesce(sqlc.narg('verified_email'), verified_email),
    verified_phone = coalesce(sqlc.narg('verified_phone'), verified_phone),
    hashed_password = coalesce(sqlc.narg('hashed_password'), hashed_password),
    password_changed_at = coalesce(sqlc.narg('password_changed_at'), password_changed_at),
    updated_at = sqlc.arg('updated_at')
WHERE
    user_id = sqlc.arg('id')
RETURNING user_id, email, username, fullname, role, verified_email, verified_phone, created_at, updated_at;


-- name: DeleteUser :exec
DELETE FROM
    users
WHERE
    user_id = $1;

-- name: CountUsers :one
SELECT
    count(*)
FROM
    users;