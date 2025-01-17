-- name: CreateAttribute :one
INSERT INTO attributes (
    attribute_name
) VALUES (
    $1
) RETURNING *;

-- name: AddBulkAttributes :copyfrom
INSERT INTO attributes (attribute_name) VALUES ($1);

-- name: GetAttributeByID :one
SELECT
    *
FROM
    attributes
WHERE
    attribute_id = $1;

-- name: GetAttributeDetailsByID :many
SELECT
    *
FROM
    attributes
JOIN
    attribute_values ON attributes.attribute_id = attribute_values.attribute_id
WHERE
    attributes.attribute_id = $1
ORDER BY
    attributes.attribute_id;

-- name: GetAttributes :many
SELECT
    *
FROM
    attributes
JOIN
    attribute_values ON attributes.attribute_id = attribute_values.attribute_id
ORDER BY
    attributes.attribute_id;

-- name: GetAttributeByName :one
SELECT
    *
FROM
    attributes
WHERE
    attribute_name = $1
LIMIT 1;

-- name: CountAttributes :one
SELECT COUNT(*) FROM attributes;

-- name: DeleteAttribute :exec
DELETE FROM
    attributes
WHERE
    attribute_id = $1;

-- name: UpdateAttribute :one
UPDATE
    attributes
SET
    attribute_name = $2
WHERE
    attribute_id = $1
RETURNING *;

------ Attribute Values ------

-- name: CreateAttributeValue :one
INSERT INTO attribute_values (
    attribute_id,
    attribute_value,
    color
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateBulkAttributeValues :copyfrom
INSERT INTO attribute_values (attribute_id, attribute_value, color) VALUES ($1, $2, $3);

-- name: GetAttributeValueByID :one
SELECT
    *
FROM
    attribute_values
WHERE
    attribute_value_id = $1 AND attribute_id = $2
LIMIT 1;

-- name: GetAttributeValueByValue :one
SELECT
    *
FROM
    attribute_values
WHERE
    attribute_value = $1
LIMIT 1;

-- name: GetAttributeValues :many
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

-- name: UpdateAttributeValue :one
UPDATE
    attribute_values
SET
    attribute_id = $2,
    attribute_value = $3,
    color = COALESCE($4, color)
WHERE
    attribute_value_id = $1
RETURNING *;

-- name: DeleteAttributeValue :exec
DELETE FROM
    attribute_values
WHERE
    attribute_value_id = $1;
