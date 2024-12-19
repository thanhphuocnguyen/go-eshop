-- name: CreateAddress :one
INSERT INTO
    user_addresses (
        user_id,
        phone,
        address_1,
        address_2,
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
        $6,
        $7
    ) RETURNING *;

-- name: GetAddress :one
SELECT
    *
FROM
    user_addresses
WHERE
    id = $1 AND user_id = $2 AND is_primary = COALESCE(sqlc.narg('is_primary'), is_primary)
LIMIT 1;

-- name: ListAddresses :many
SELECT
    *
FROM
    user_addresses
WHERE
    user_id = $1
ORDER BY
    id
LIMIT $2
OFFSET $3;

-- name: UpdateAddress :one
UPDATE
    user_addresses
SET
    phone = coalesce(sqlc.narg('phone'), phone),
    address_1 = coalesce(sqlc.narg('address_1'), address_1),
    address_2 = coalesce(sqlc.narg('address_2'), address_2),
    ward = coalesce(sqlc.narg('ward'), ward),
    district = coalesce(sqlc.narg('district'), district),
    city = coalesce(sqlc.narg('city'), city),
    is_primary = coalesce(sqlc.narg('is_primary'), is_primary),
    updated_at = now()
WHERE
    id = sqlc.arg('id') AND user_id = sqlc.arg('user_id')
RETURNING *;

-- name: DeleteAddress :exec
DELETE FROM
    user_addresses
WHERE
    id = $1 AND user_id = $2;

-- name: SetPrimaryAddress :exec
UPDATE
    user_addresses
SET
    is_primary = $1,
    updated_at = now()
WHERE
    user_id = $2 AND id = $3;
