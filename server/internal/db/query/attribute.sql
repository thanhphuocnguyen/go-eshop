-- name: CreateAttribute :one
INSERT INTO attributes (
    attribute_name
) VALUES (
    $1
) RETURNING *;

-- name: CreateAttributeValue :one
INSERT INTO attribute_values (
    attribute_id,
    attribute_value
) VALUES (
    $1, $2
) RETURNING *;

-- name: DeleteAttribute :exec
DELETE FROM
    attributes
WHERE
    attribute_id = $1;

-- name: DeleteAttributeValue :exec
DELETE FROM
    attribute_values
WHERE
    attribute_value_id = $1;

-- name: GetAttributeByID :one
SELECT
    *
FROM
    attributes
WHERE
    attribute_id = $1
LIMIT 1;

-- name: GetAttributeValueByID :one
SELECT
    *
FROM
    attribute_values
WHERE
    attribute_value_id = $1
LIMIT 1;

-- name: ListAttributes :many
SELECT
    *
FROM
    attributes
ORDER BY
    attribute_id
LIMIT $1
OFFSET $2;

-- name: ListAttributeValues :many
SELECT
    *
FROM
    attribute_values
WHERE
    attribute_id = $1
ORDER BY
    attribute_value_id
LIMIT $2
OFFSET $3;

-- name: ListVariantAttributes :many
SELECT
    *
FROM
    variant_attributes
WHERE
    variant_id = $1
ORDER BY
    variant_attribute_id
LIMIT $2
OFFSET $3;

-- name: UpdateAttribute :one
UPDATE
    attributes
SET
    attribute_name = $2
WHERE
    attribute_id = $1
RETURNING *;

-- name: UpdateAttributeValue :one
UPDATE
    attribute_values
SET
    attribute_id = $2,
    attribute_value = $3
WHERE
    attribute_value_id = $1
RETURNING *;


-- name: GetAttributeByName :one
SELECT
    *
FROM
    attributes
WHERE
    attribute_name = $1
LIMIT 1;

-- name: GetAttributeValueByValue :one
SELECT
    *
FROM
    attribute_values
WHERE
    attribute_value = $1
LIMIT 1;

-- name: SeedAttributes :copyfrom
INSERT INTO attributes (attribute_name) VALUES ($1);

-- name: SeedAttributeValues :copyfrom
INSERT INTO attribute_values (attribute_id, attribute_value) VALUES ($1, $2);

-- name: CountAttributes :one
SELECT COUNT(*) FROM attributes;
