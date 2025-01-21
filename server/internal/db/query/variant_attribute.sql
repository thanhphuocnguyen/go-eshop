-- name: CreateVariantAttribute :one
INSERT INTO variant_attributes (
    variant_id,
    attribute_id,
    value
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetVariantAttribute :one
SELECT
    *
FROM
    variant_attributes
WHERE
    variant_attribute_id = $1
LIMIT 1;

-- name: GetVariantAttributes :many
SELECT
    *
FROM
    variant_attributes
WHERE
    variant_id = $1;

-- name: UpdateVariantAttribute :one
UPDATE
    variant_attributes
SET
    value = COALESCE(sqlc.narg('value'), value)
WHERE
    variant_attribute_id = $1
RETURNING *;

-- name: RemoveVariantAttribute :exec
DELETE FROM
    variant_attributes
WHERE
    variant_attribute_id = $1;

-- name: CreateBulkVariantAttributes :copyfrom
INSERT INTO variant_attributes (
    variant_id,
    attribute_id,
    value
) VALUES
    ($1, $2, $3);