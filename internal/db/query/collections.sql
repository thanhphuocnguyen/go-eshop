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
    p.base_sku as product_sku, p.slug as product_slug,
    pi.id as image_id, pi.image_url
FROM collections AS c
LEFT JOIN products AS p ON c.id = p.collection_id
LEFT JOIN product_images AS pi ON p.id = pi.product_id
WHERE c.id = ANY(sqlc.narg('ids')::UUID[])
GROUP BY c.id, p.id, pi.id, pi.image_url
LIMIT $1 OFFSET $2;

-- name: GetCollections :many
SELECT * FROM collections WHERE  published = COALESCE(sqlc.narg('published'), published) ORDER BY display_order LIMIT $1 OFFSET $2;

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
