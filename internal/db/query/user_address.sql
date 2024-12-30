-- name: CreateAddress :one
INSERT INTO
    user_addresses (
        user_id,
        phone,
        street,
        ward,
        district,
        city
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    ) RETURNING *;

-- name: GetAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    id = $1 AND user_id = $2 AND is_deleted = false
LIMIT 1;

-- name: GetAddresses :many
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = $1 AND is_deleted = false
ORDER BY
    is_primary DESC, id ASC;

-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone = coalesce(sqlc.narg('phone'), phone),
    street = coalesce(sqlc.narg('street'), street),
    ward = coalesce(sqlc.narg('ward'), ward),
    district = coalesce(sqlc.narg('district'), district),
    city = coalesce(sqlc.narg('city'), city),
    is_primary = coalesce(sqlc.narg('is_primary'), is_primary)
WHERE
    id = sqlc.arg('id') AND user_id = sqlc.arg('user_id') AND is_deleted = false
RETURNING *;

-- name: DeleteAddress :exec
UPDATE
    user_addresses
SET
    is_deleted = true
WHERE
    id = $1 AND user_id = $2;

-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    is_primary = $1
WHERE
    user_id = $2 AND id = $3 AND is_deleted = false;

-- name: ResetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    is_primary = false
WHERE
    user_id = $1 AND is_primary = true;

-- name: GetPrimaryAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = $1 AND is_primary = true AND is_deleted = false
LIMIT 1;

-- name: CountAddresses :one
SELECT
    COUNT(*)
FROM
    user_addresses;