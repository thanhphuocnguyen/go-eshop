-- name: CreateCollection :one
INSERT INTO collections (name, slug, description,  image_url, image_id) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetCollectionByID :one
SELECT * FROM collections  WHERE id = $1 LIMIT 1;

-- name: GetCollectionBySlug :one
SELECT * FROM collections WHERE slug = $1 LIMIT 1;

-- name: GetCollectionsByIDs :many
SELECT 
    c.*, 
    p.name as product_name, p.id, p.description,
    p.base_price as product_price, 
    p.base_sku as product_sku, p.slug as product_slug
FROM collections AS c
LEFT JOIN collection_products AS cp ON c.id = cp.collection_id
LEFT JOIN products AS p ON cp.product_id = p.id
WHERE c.id = ANY(sqlc.narg('ids')::UUID[])
GROUP BY c.id, p.id
LIMIT $1 OFFSET $2;

-- name: GetCollections :many
SELECT * FROM collections WHERE  published = COALESCE(sqlc.narg('published'), published) ORDER BY display_order LIMIT $1 OFFSET $2;

-- name: GetDisplayCollectionProducts :many
SELECT 
    p.*, COUNT(pv.id) AS variant_count, MIN(pv.price) AS price
FROM collections AS c
JOIN collection_products AS cp ON c.id = cp.collection_id
JOIN products AS p ON cp.product_id = p.id
LEFT JOIN product_variants AS pv ON pv.product_id = p.id
WHERE c.id = $1 AND p.is_active = TRUE
GROUP BY p.id
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateCollectionWith :one
UPDATE collections
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    image_id = COALESCE(sqlc.narg('image_id'), image_id),
    description = COALESCE(sqlc.narg('description'), description),
    slug = COALESCE(sqlc.narg('slug'), slug),
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE id = $1 RETURNING *;

-- name: DeleteCollection :exec
DELETE FROM collections WHERE id = $1;

-- name: CountCollections :one
SELECT count(*) FROM collections;

-- name: SeedCollections :copyfrom
INSERT INTO collections (name, description, image_url) VALUES ($1, $2, $3);

-- name: AddProductsToCollection :copyfrom
INSERT INTO collection_products (collection_id, product_id) VALUES ($1, $2);

-- name: RemoveProductsFromCollection :exec
DELETE FROM collection_products WHERE product_id = $1;