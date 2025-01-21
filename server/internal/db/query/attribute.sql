-- name: CreateAttribute :one
INSERT INTO attributes (
    name
) VALUES (
    $1
) RETURNING *;

-- name: AddBulkAttributes :copyfrom
INSERT INTO attributes (name) VALUES ($1);

-- name: GetAttributeByID :one
SELECT
    *
FROM
    attributes
WHERE
    attribute_id = $1;

-- name: GetAttributes :many
SELECT
    *
FROM
    attributes
ORDER BY
    attributes.attribute_id;

-- name: GetAttributeByName :one
SELECT
    *
FROM
    attributes
WHERE
    name = $1
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
    name = $2
WHERE
    attribute_id = $1
RETURNING *;
