-- name: CreateBrand :one
INSERT INTO brands (name, slug, description, image_url, image_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetBrandByID :one
SELECT * FROM brands WHERE id = $1 LIMIT 1;

-- name: GetBrandBySlug :one
SELECT * FROM brands WHERE slug = $1 LIMIT 1;

-- name: GetBrandsByIDs :many
SELECT c.*, p.name as product_name, p.id as product_id, p.description, p.base_price as product_price, p.base_sku as product_sku, p.slug as product_slug
FROM brands AS c
LEFT JOIN products AS p ON p.brand_id = c.id
WHERE c.id = ANY(sqlc.narg('ids')::UUID[])
GROUP BY c.id, p.id
LIMIT $1 OFFSET $2;

-- name: GetBrands :many
SELECT * FROM brands WHERE published = COALESCE(sqlc.narg('published'), published) ORDER BY display_order LIMIT $1 OFFSET $2;

-- name: UpdateBrandWith :one
UPDATE brands
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url),
    image_id = COALESCE(sqlc.narg('image_id'), image_id),
    description = COALESCE(sqlc.narg('description'), description),
    slug = COALESCE(sqlc.narg('slug'), slug),
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brands WHERE id = $1;

-- name: CountBrands :one
SELECT count(*) FROM brands;

-- name: SeedBrands :copyfrom
INSERT INTO brands (name, slug, description, display_order, published) VALUES ($1, $2, $3, $4, $5);