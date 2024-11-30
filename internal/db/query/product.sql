-- name: CreateProduct :one
INSERT INTO
    products (
        name,
        description,
        sku,
        image_url,
        stock,
        price
    )
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    )
RETURNING *;

-- name: GetProduct :one
SELECT
    *
FROM
    products
WHERE
    id = $1
LIMIT 1;

-- name: ListProducts :many
SELECT
    *
FROM
    products
ORDER BY
    id
LIMIT $1
OFFSET $2;


-- name: UpdateProduct :one
UPDATE
    products
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    sku = coalesce(sqlc.narg('sku'), sku),
    image_url = coalesce(sqlc.narg('image_url'), image_url),
    stock = coalesce(sqlc.narg('stock'), stock),
    price = coalesce(sqlc.narg('price'), price),
    updated_at = sqlc.arg('updated_at')
WHERE
    id = sqlc.arg('id')
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM
    products
WHERE
    id = $1;

-- name: ArchiveProduct :exec
UPDATE
    products
SET
    archived = true
WHERE
    id = $1
RETURNING *;