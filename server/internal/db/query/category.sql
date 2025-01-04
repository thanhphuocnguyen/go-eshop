-- name: CreateCollection :one
INSERT INTO categories (name, sort_order, image_url, published)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCollection :one
SELECT sqlc.embed(categories), sqlc.embed(images), sqlc.embed(products), sqlc.embed(category_products)
FROM categories
JOIN category_products ON categories.category_id = category
JOIN products ON category_products.product_id = products.product_id
JOIN images ON products.product_id = images.product_id AND images.primary = TRUE
WHERE categories.category_id = $1
LIMIT 1;

-- name: GetCollectionByName :one
SELECT *
FROM categories
WHERE name = $1
LIMIT 1;

-- name: GetCollections :many
SELECT *
FROM categories
WHERE categories.id = ANY(bigint[] $1);
ORDER BY sort_order;
