-- name: CreateCategory :one
INSERT INTO categories (name, image_url, sort_order)
VALUES (
    $1,
    $2,
    COALESCE(sqlc.narg('sort_order'), COALESCE((SELECT MAX(sort_order) + 1 FROM categories), 1))
)
RETURNING *;

-- name: GetCategoryByID :one
SELECT c.* FROM categories c WHERE c.category_id = $1 LIMIT 1;

-- name: GetCategoriesByIDs :many
SELECT 
    c.*, 
    p.name as product_name, p.product_id, p.description,
    MIN(pv.price) as price_from, 
    MAX(pv.price) as price_to, 
    MAX(pv.discount) as discount, 
    MIN(pv.stock_quantity) as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM categories AS c
LEFT JOIN products AS p ON c.category_id = p.category_id
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE c.category_id = ANY(sqlc.narg('category_ids')::int[]) AND c.published = sqlc.narg('published')
GROUP BY c.category_id, p.product_id, img.image_id, img.image_url
ORDER BY c.sort_order;

-- name: GetCategories :many
SELECT 
    c.*, 
    p.name as product_name, p.product_id, p.description,
    MIN(pv.price) as price_from, 
    MAX(pv.price) as price_to, 
    MAX(pv.discount) as discount, 
    MIN(pv.stock_quantity) as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM categories AS c
LEFT JOIN products AS p ON c.category_id = p.category_id
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE c.published = sqlc.narg('published')
GROUP BY c.category_id, p.product_id, img.image_id, img.image_url
ORDER BY c.sort_order;

-- name: UpdateCategoryWith :one
UPDATE categories
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order), 
    published = COALESCE(sqlc.narg('published'), published)
WHERE category_id = $1
RETURNING *;


-- name: DeleteCategory :exec
DELETE FROM categories WHERE category_id = $1;

-- name: CountCategories :one
SELECT count(*)
FROM categories
WHERE category_id = COALESCE(sqlc.narg('category_id'), category_id);

-- name: SeedCategories :copyfrom
INSERT INTO categories (name, image_url, sort_order, published) VALUES ($1, $2, $3, $4);

-- name: GetCategoryMaxSortOrder :one
SELECT COALESCE(MAX(sort_order)::smallint, 0) AS max_sort_order FROM categories;