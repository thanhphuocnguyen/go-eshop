-- name: CreateCollection :one
INSERT INTO categories (name, description, sort_order, image_url, published)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCollection :many
SELECT 
    c.*, 
    p.name, p.description, p.price, p.discount, 
    cp.product_id, 
    i.image_id, i.image_url
FROM categories c
JOIN category_products cp ON cp.category_id = c.category_id
JOIN products p ON cp.product_id = p.product_id AND p.published = TRUE
LEFT JOIN images i ON p.product_id = i.product_id AND i.primary = TRUE
WHERE categories.category_id = $1;

-- name: GetCollectionByName :one
SELECT *
FROM categories c
JOIN category_products cp ON c.category_id = cp.category_id
JOIN products p ON cp.product_id = p.product_id AND p.published = TRUE
WHERE c.name = $1 AND c.published = TRUE
LIMIT 1;

-- name: GetCollections :many
SELECT 
    c.category_id, c.name, c.description, c.sort_order, c.image_url, c.published,
    p.name, p.description, p.price, p.discount, 
    cp.product_id, 
    images.image_id, images.image_url
FROM categories c
JOIN category_products cp ON c.category_id = cp.category_id
JOIN products p ON category_products.product_id = p.product_id AND p.published = TRUE
LEFT JOIN images ON p.product_id = images.product_id AND images.primary = TRUE
WHERE categories.category_id = ANY(sqlc.arg('category_ids')::int[]) AND published = TRUE
ORDER BY c.sort_order, cp.sort_order;

-- name: UpdateCollection :one
UPDATE categories
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    description = COALESCE(sqlc.narg('description'), description), 
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE category_id = $1
RETURNING *;


-- name: RemoveCollection :exec
DELETE FROM categories
WHERE category_id = $1;