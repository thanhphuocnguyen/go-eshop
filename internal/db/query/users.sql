-- name: CreateUser :one
INSERT INTO users (email, username, phone_number, first_name, last_name, hashed_password, role_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserDetailsByID :one
SELECT u.id, u.email, u.username, u.first_name, u.last_name, 
    u.phone_number, u.role_id, u.verified_email, 
    u.verified_phone, u.created_at, u.updated_at, 
    ur.code AS role_code, 
    COUNT(CASE WHEN ord.status IN ('completed', 'delivered') AND p.status = 'success' THEN 1 END) AS total_orders,
    COALESCE(SUM(CASE 
        WHEN ord.status IN ('completed', 'delivered') AND p.status = 'success'
        THEN ord.total_price - COALESCE(du.total_discount, 0)
        ELSE 0 
    END), 0) AS total_spent,
    COALESCE(MAX(CASE 
        WHEN ord.status IN ('completed', 'delivered') AND p.status = 'success'
        THEN ord.total_price - COALESCE(du.total_discount, 0)
        ELSE 0 
    END), 0) AS largest_order_amount
FROM users u
JOIN user_roles ur ON u.role_id = ur.id
LEFT JOIN orders ord ON u.id = ord.user_id
LEFT JOIN payments p ON ord.id = p.order_id
LEFT JOIN (
    SELECT order_id, SUM(discount_amount) as total_discount
    FROM discount_usage
    GROUP BY order_id
) du ON ord.id = du.order_id
WHERE u.id = $1
GROUP BY u.id, ur.code
LIMIT 1;

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
INSERT INTO user_addresses (user_id, phone_number, street, ward, district, city, "is_default") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetAddress :one
SELECT * FROM user_addresses WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: DeleteAddress :exec
DELETE FROM user_addresses WHERE id = $1 AND user_id = $2;

-- name: GetDefaultAddress :one
SELECT * FROM user_addresses WHERE user_id = $1 AND "is_default" = TRUE LIMIT 1;

-- name: GetAddresses :many
SELECT * FROM user_addresses WHERE user_id = $1 ORDER BY "is_default" DESC, id ASC;

-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone_number = coalesce(sqlc.narg('phone_number'), phone_number),
    street = coalesce(sqlc.narg('street'), street),
    ward = coalesce(sqlc.narg('ward'), ward),
    district = coalesce(sqlc.narg('district'), district),
    city = coalesce(sqlc.narg('city'), city),
    "is_default" = coalesce(sqlc.narg('is_default'), "is_default")
WHERE
    id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;


-- name: SetPrimaryAddress :exec
UPDATE user_addresses SET "is_default" = $1 WHERE id = $2 AND user_id = $3;

-- name: ResetPrimaryAddress :exec
UPDATE user_addresses SET "is_default" = FALSE WHERE user_id = $1 AND "is_default" = TRUE;

-- name: CountAddresses :one
SELECT COUNT(*) FROM user_addresses;

-- name: SeedAddresses :copyfrom
INSERT INTO user_addresses (user_id, phone_number, street, ward, district, city, "is_default") VALUES ($1,$2,$3,$4,$5,$6,$7);

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

-- name: GetUserTotalSpent :one
SELECT COALESCE(SUM(o.total_price - COALESCE(du.total_discount, 0)), 0) as total_spent
FROM orders o
JOIN payments p ON o.id = p.order_id
LEFT JOIN (
    SELECT order_id, SUM(discount_amount) as total_discount
    FROM discount_usage
    GROUP BY order_id
) du ON o.id = du.order_id
WHERE o.user_id = $1 
  AND o.status IN ('completed', 'delivered')
  AND p.status = 'success';