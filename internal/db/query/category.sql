-- name: CreateCollection :one
INSERT INTO categories (name, sort_order, image_url, published)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCollection :one
SELECT *
FROM categories
WHERE id = $1
LIMIT 1;
