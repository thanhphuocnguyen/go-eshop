-- name: CreateUser :one
INSERT INTO
    users (id,email,username,phone,fullname,hashed_password,role)
VALUES
    ($1,$2,$3,$4,$5,$6,$7)
RETURNING id, email, username, fullname, role, verified_email, verified_phone, created_at, updated_at;

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
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = coalesce(sqlc.narg('email'), email),
    fullname = coalesce(sqlc.narg('fullname'), fullname),
    role = coalesce(sqlc.narg('role'), role),
    phone = coalesce(sqlc.narg('phone'), phone),
    verified_email = coalesce(sqlc.narg('verified_email'), verified_email),
    verified_phone = coalesce(sqlc.narg('verified_phone'), verified_phone),
    hashed_password = coalesce(sqlc.narg('hashed_password'), hashed_password),
    password_changed_at = coalesce(sqlc.narg('password_changed_at'), password_changed_at),
    updated_at = sqlc.arg('updated_at')
WHERE
    id = sqlc.arg('id')
RETURNING id, email, username, fullname, role, verified_email, verified_phone, created_at, updated_at;


-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CountUsers :one
SELECT
    count(*)
FROM
    users;

-- name: SeedUsers :copyfrom
INSERT INTO
    users (id,email,username,phone,fullname,hashed_password,role)
VALUES
    ($1, $2, $3, $4, $5, $6, $7);

-- User Address Queries
-- name: CreateAddress :one
INSERT INTO 
    user_addresses (user_id, phone, street, ward, district, city, "default")
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    id = $1 AND user_id = $2 AND deleted = FALSE
LIMIT 1;


-- name: GetDefaultAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = $1 AND deleted = FALSE AND "default" = TRUE
LIMIT 1;


-- name: GetAddresses :many
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = $1 AND deleted = FALSE
ORDER BY
    "default" DESC, id ASC;

-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone = coalesce(sqlc.narg('phone'), phone),
    street = coalesce(sqlc.narg('street'), street),
    ward = coalesce(sqlc.narg('ward'), ward),
    district = coalesce(sqlc.narg('district'), district),
    city = coalesce(sqlc.narg('city'), city),
    "default" = coalesce(sqlc.narg('default'), "default")
WHERE
    id = sqlc.arg('id') AND user_id = sqlc.arg('user_id') AND deleted = FALSE
RETURNING *;

-- name: DeleteAddress :exec
UPDATE
    user_addresses
SET
    deleted = TRUE,
    updated_at = now()
WHERE
    id = $1 AND user_id = $2;

-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    "default" = $1
WHERE
    id = $2 AND user_id = $3 AND deleted = FALSE;

-- name: ResetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    "default" = FALSE
WHERE
    user_id = $1 AND "default" = TRUE;

-- name: CountAddresses :one
SELECT
    COUNT(*)
FROM
    user_addresses;

-- name: SeedAddresses :copyfrom
INSERT INTO
    user_addresses (
        user_id,
        phone,
        street,
        ward,
        district,
        city,
        "default"
    )
VALUES
    ($1,$2,$3,$4,$5,$6,$7);

-- Verification Token Queries
-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (id, email, verify_code) VALUES ($1, $2, $3) RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails WHERE id = $1 AND email = $2;

-- name: GetVerifyEmailByID :one
SELECT * FROM verify_emails WHERE id = $1;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = $1 AND verify_code = $2 AND expired_at > now()
RETURNING *;