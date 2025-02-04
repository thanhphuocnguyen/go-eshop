-- name: CreateCategory :one
INSERT INTO categories (name, description, sort_order)
VALUES (
    $1,
    $2,
    COALESCE(sqlc.narg('sort_order'), COALESCE((SELECT MAX(sort_order) + 1 FROM categories), 1))
)
RETURNING *;

-- name: GetCategoryByID :many
SELECT c.* FROM categories c WHERE c.category_id = $1;

-- name: GetCategoryWithProduct :many
SELECT 
    c.*, 
    p.name as product_name, p.description,
    cp.product_id,
    MIN(pv.price)::decimal as price_from, 
    MAX(pv.price)::decimal as price_to, 
    MAX(pv.discount)::smallint as discount, 
    MIN(pv.stock_quantity)::smallint as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM categories AS c
JOIN category_products AS cp ON cp.category_id = c.category_id
JOIN products AS p ON cp.product_id = p.product_id
JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE c.category_id = $1
GROUP BY c.category_id, p.product_id, img.image_id, img.image_url, cp.product_id;

-- name: GetCategoryByName :one
SELECT *
FROM categories c
JOIN category_products cp ON c.category_id = cp.category_id
JOIN products p ON cp.product_id = p.product_id
WHERE c.name = $1 AND c.published = TRUE
LIMIT 1;

-- name: GetCategories :many
SELECT 
    categories.*,
    COUNT(category_products.product_id) as product_count
FROM 
    categories
LEFT JOIN category_products 
    ON categories.category_id = category_products.category_id
WHERE 
    categories.published = COALESCE(sqlc.narg('published'), categories.published)
GROUP BY categories.category_id
ORDER BY categories.sort_order;
-- name: GetCategoriesInIDs :many
SELECT 
    *
FROM categories
WHERE 
    category_id = ANY(sqlc.arg(id_list)::int[])
    AND published = TRUE
ORDER BY sort_order;

-- name: UpdateCategoryWith :one
UPDATE categories
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    description = COALESCE(sqlc.narg('description'), description), 
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order), 
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE category_id = $1
RETURNING *;


-- name: RemoveCategory :exec
DELETE FROM categories WHERE category_id = $1;

-- name: CountCategories :one
SELECT count(*)
FROM categories
WHERE category_id = COALESCE(sqlc.narg('category_id'), category_id);

-- name: SeedCategories :copyfrom
INSERT INTO categories (name, description, sort_order, published) VALUES ($1, $2, $3, $4);

-- name: GetCategoryMaxSortOrder :one
SELECT COALESCE(MAX(sort_order)::smallint, 0) AS max_sort_order
FROM category_products;

-- Category Products

-- name: AddProductToCategory :one
INSERT INTO
    category_products (category_id, product_id, sort_order)
VALUES
    ($1, $2, $3)
RETURNING *;

-- name: RemoveProductFromCategory :exec
DELETE FROM
    category_products
WHERE
    category_id = $1
    AND product_id = $2;

-- name: GetCategoryProduct :one
SELECT
    p.*
FROM
    products p
JOIN category_products cp ON p.product_id = cp.product_id
WHERE
    cp.category_id = $1
    AND cp.product_id = $2;

-- name: GetCategoryProducts :many
SELECT
    p.*
FROM
    products p
    JOIN category_products cp ON p.product_id = cp.product_id
WHERE
    cp.category_id = $1;

-- name: GetMaxSortOrderInCategory :one
SELECT
    max(sort_order)
FROM
    category_products
WHERE
    category_id = $1;

-- name: UpdateProductSortOrderInCategory :exec
UPDATE
    category_products
SET
    sort_order = $3
WHERE
    category_id = $1
    AND product_id = $2;