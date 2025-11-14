-- name: CreateUser :one
INSERT INTO users (email, username, phone_number, first_name, last_name, hashed_password, role_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE
    users
SET
    email = coalesce(sqlc.narg('email'), email),
    first_name = coalesce(sqlc.narg('first_name'), first_name),
    last_name = coalesce(sqlc.narg('last_name'), last_name),
    role_id = coalesce(sqlc.narg('role_id'), role_id),
    phone_number = coalesce(sqlc.narg('phone_number'), phone_number),
    verified_email = coalesce(sqlc.narg('verified_email'), verified_email),
    verified_phone = coalesce(sqlc.narg('verified_phone'), verified_phone),
    hashed_password = coalesce(sqlc.narg('hashed_password'), hashed_password),
    password_changed_at = coalesce(sqlc.narg('password_changed_at'), password_changed_at),
    updated_at = sqlc.arg('updated_at')
WHERE
    id = sqlc.arg('id')
RETURNING id, email, username, first_name, last_name, role_id, verified_email, verified_phone, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: SeedUsers :copyfrom
INSERT INTO users (email, username, phone_number, first_name, last_name, hashed_password, role_id) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- User Address Queries
-- name: CreateAddress :one
INSERT INTO user_addresses (user_id, phone_number, street, ward, district, city, "default") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetAddress :one
SELECT * FROM user_addresses WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: DeleteAddress :exec
DELETE FROM user_addresses WHERE id = $1 AND user_id = $2;

-- name: GetDefaultAddress :one
SELECT * FROM user_addresses WHERE user_id = $1 AND "default" = TRUE LIMIT 1;

-- name: GetAddresses :many
SELECT * FROM user_addresses WHERE user_id = $1 ORDER BY "default" DESC, id ASC;

-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone_number = coalesce(sqlc.narg('phone_number'), phone_number),
    street = coalesce(sqlc.narg('street'), street),
    ward = coalesce(sqlc.narg('ward'), ward),
    district = coalesce(sqlc.narg('district'), district),
    city = coalesce(sqlc.narg('city'), city),
    "default" = coalesce(sqlc.narg('default'), "default")
WHERE
    id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;


-- name: SetPrimaryAddress :exec
UPDATE user_addresses SET "default" = $1 WHERE id = $2 AND user_id = $3;

-- name: ResetPrimaryAddress :exec
UPDATE user_addresses SET "default" = FALSE WHERE user_id = $1 AND "default" = TRUE;

-- name: CountAddresses :one
SELECT COUNT(*) FROM user_addresses;

-- name: SeedAddresses :copyfrom
INSERT INTO user_addresses (user_id, phone_number, street, ward, district, city, "default") VALUES ($1,$2,$3,$4,$5,$6,$7);

-- Verification Token Queries
-- name: CreateVerifyEmail :one
INSERT INTO email_verifications (user_id, email, verify_code) VALUES ($1, $2, $3) RETURNING *;

-- name: GetVerifyEmailByVerifyCode :one
SELECT * FROM email_verifications WHERE verify_code = $1 AND expired_at > now() AND is_used = FALSE;

-- name: GetVerifyEmailByID :one
SELECT * FROM email_verifications WHERE id = $1;

-- name: UpdateVerifyEmail :one
UPDATE email_verifications SET is_used = TRUE WHERE id = $1 AND verify_code = $2 AND expired_at > now() RETURNING *;

-- Roles Queries
-- name: GetRoleByCode :one
SELECT * FROM user_roles WHERE code = $1 LIMIT 1;

-- name: GetRoleByID :one
SELECT * FROM user_roles WHERE id = $1 LIMIT 1;