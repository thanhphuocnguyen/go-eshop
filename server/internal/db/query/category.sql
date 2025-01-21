-- name: CreateCollection :one
INSERT INTO categories (name, description, sort_order)
VALUES (
    $1,
    $2,
    COALESCE(sqlc.narg('sort_order'), COALESCE((SELECT MAX(sort_order) + 1 FROM categories), 1))
)
RETURNING *;

-- name: GetCollection :many
SELECT 
    c.*, 
    p.name, p.description,
    cp.product_id,
    MIN(pv.price)::decimal as price_from, MAX(pv.price)::decimal as price_to, MAX(pv.discount)::int as discount, MIN(pv.stock_quantity)::int as stock_quantity, COUNT(pv.product_variant_id) as variant_count,
    i.image_id, i.image_url
FROM categories c
JOIN category_products cp ON cp.category_id = c.category_id
JOIN products p ON cp.product_id = p.product_id AND p.published = TRUE
JOIN product_variants pv ON p.product_id = pv.product_id
LEFT JOIN images i ON p.product_id = i.product_id AND i.primary = TRUE
WHERE categories.category_id = $1
GROUP BY c.category_id, p.product_id, i.image_id;

-- name: GetCollectionByName :one
SELECT *
FROM categories c
JOIN category_products cp ON c.category_id = cp.category_id
JOIN products p ON cp.product_id = p.product_id AND p.published = TRUE
WHERE c.name = $1 AND c.published = TRUE
LIMIT 1;

-- name: GetCollections :many
SELECT 
    *
FROM categories
WHERE 
    published = TRUE
ORDER BY sort_order;
-- name: GetCollectionsInIDs :many
SELECT 
    *
FROM categories
WHERE 
    category_id = ANY(sqlc.arg(id_list)::int[])
    AND published = TRUE
ORDER BY sort_order;

-- name: UpdateCollection :one
UPDATE categories
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    description = COALESCE(sqlc.narg('description'), description), 
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order), 
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE category_id = $1
RETURNING *;


-- name: RemoveCollection :exec
DELETE FROM categories
WHERE category_id = $1;

-- name: CountCollections :one
SELECT count(*)
FROM categories
WHERE category_id = COALESCE(sqlc.narg('category_id'), category_id);

-- name: SeedCollections :copyfrom
INSERT INTO categories (name, description, sort_order, published) VALUES ($1, $2, $3, $4);

-- name: GetCollectionMaxSortOrder :one
SELECT COALESCE(MAX(sort_order), 0)::smallint AS max_sort_order
FROM category_products;