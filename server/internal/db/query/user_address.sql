-- name: CreateAddress :one
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
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    ) RETURNING *;

-- name: GetAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    user_address_id = $1 AND user_id = $2 AND deleted = FALSE
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
    "default" DESC, user_address_id ASC;

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
    user_address_id = sqlc.arg('user_address_id') AND user_id = sqlc.arg('user_id') AND deleted = FALSE
RETURNING *;

-- name: DeleteAddress :exec
UPDATE
    user_addresses
SET
    deleted = TRUE,
    updated_at = now()
WHERE
    user_address_id = $1 AND user_id = $2;

-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    "default" = $1
WHERE
    user_id = $2 AND user_address_id = $3 AND deleted = FALSE;

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