-- name: CreateCollection :one
INSERT INTO collections 
    (id, name, slug, description, image_url, image_id)
VALUES 
    ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCollectionByID :one
SELECT c.*
FROM collections c 
WHERE c.id = $1
LIMIT 1;

-- name: GetCollectionByIDWithProducts :many
SELECT 
    c.*,
    p.name as product_name, p.id as product_id, p.description as product_description, 
    p.base_price as product_price, 
    p.base_sku as product_sku, p.slug as product_slug,
    img.id as image_id, img.url as product_image_url
FROM 
    collections c
LEFT JOIN
    products p ON c.id = p.collection_id
LEFT JOIN
    image_assignments ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN
    images img ON img.id = ia.image_id
WHERE c.id = $1;

-- name: GetCollectionsByIDs :many
SELECT 
    c.*, 
    p.name as product_name, p.id, p.description,
    p.base_price as product_price, 
    p.base_sku as product_sku, p.slug as product_slug,
    img.id as image_id, img.url
FROM collections AS c
LEFT JOIN products AS p ON c.id = p.collection_id
LEFT JOIN image_assignments AS ia ON p.id = ia.entity_id AND ia.entity_type = 'product'
LEFT JOIN images AS img ON img.id = ia.image_id
WHERE c.id = ANY(sqlc.narg('ids')::int[])
GROUP BY c.id, p.id, img.id, img.url
LIMIT $1 OFFSET $2;

-- name: GetCollections :many
SELECT 
    c.*
FROM collections AS c
LIMIT $1 OFFSET $2;

-- name: UpdateCollectionWith :one
UPDATE collections
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    image_id = COALESCE(sqlc.narg('image_id'), image_id),
    description = COALESCE(sqlc.narg('description'), description),
    remarkable = COALESCE(sqlc.narg('remarkable'), remarkable),
    slug = COALESCE(sqlc.narg('slug'), slug),
    published = COALESCE(sqlc.narg('published'), published),
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeleteCollection :exec
DELETE FROM collections WHERE id = $1;

-- name: CountCollections :one
SELECT count(*)
FROM collections;

-- name: SeedCollections :copyfrom
INSERT INTO collections (name, description, image_url) VALUES ($1, $2, $3);
