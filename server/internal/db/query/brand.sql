-- name: CreateBrand :one
INSERT INTO brands (name, description, image_url)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetBrandByID :one
SELECT c.* FROM brands c WHERE c.brand_id = $1 LIMIT 1;

-- name: GetBrandsByIDs :many
SELECT 
    c.*, 
    p.name as product_name, p.product_id, p.description,
    MIN(pv.price) as price_from, 
    MAX(pv.price) as price_to, 
    MAX(pv.discount) as discount, 
    MIN(pv.stock_quantity) as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM brands AS c
LEFT JOIN products AS p ON p.brand_id = c.brand_id
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE c.brand_id = ANY(sqlc.narg('brand_ids')::int[])
GROUP BY c.brand_id, p.product_id, img.image_id, img.image_url
LIMIT $1 OFFSET $2;

-- name: GetBrands :many
SELECT 
    c.*, 
    p.name as product_name, p.product_id, p.description,
    MIN(pv.price) as price_from, 
    MAX(pv.price) as price_to, 
    MAX(pv.discount) as discount, 
    MIN(pv.stock_quantity) as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM brands AS c
LEFT JOIN products AS p ON p.brand_id = c.brand_id
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
GROUP BY c.brand_id, p.product_id, img.image_id, img.image_url
LIMIT $1 OFFSET $2;

-- name: UpdateBrandWith :one
UPDATE brands
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = now()
WHERE brand_id = $1
RETURNING *;


-- name: DeleteBrand :exec
DELETE FROM brands WHERE brand_id = $1;

-- name: CountBrands :one
SELECT count(*)
FROM brands
WHERE brand_id = COALESCE(sqlc.narg('brand_id'), brand_id);