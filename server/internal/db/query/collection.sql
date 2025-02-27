-- name: CreateCollection :one
INSERT INTO collections (name, description, image_url)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetCollectionByID :one
SELECT c.* FROM collections c WHERE c.collection_id = $1 LIMIT 1;

-- name: GetCollectionsByIDs :many
SELECT 
    c.*, 
    p.name as product_name, p.product_id, p.description,
    MIN(pv.price) as price_from, 
    MAX(pv.price) as price_to, 
    MAX(pv.discount) as discount, 
    MIN(pv.stock_quantity) as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM collections AS c
LEFT JOIN products AS p ON c.collection_id = p.collection_id
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
WHERE c.collection_id = ANY(sqlc.narg('collection_ids')::int[])
GROUP BY c.collection_id, p.product_id, img.image_id, img.image_url
LIMIT $1 OFFSET $2;

-- name: GetCollections :many
SELECT 
    c.*, 
    p.name as product_name, p.product_id, p.description,
    MIN(pv.price) as price_from, 
    MAX(pv.price) as price_to, 
    MAX(pv.discount) as discount, 
    MIN(pv.stock_quantity) as stock_quantity, 
    COUNT(pv.variant_id) as variant_count,
    img.image_id, img.image_url
FROM collections AS c
LEFT JOIN products AS p ON c.collection_id = p.collection_id
LEFT JOIN product_variants AS pv ON p.product_id = pv.product_id
LEFT JOIN images AS img ON p.product_id = img.product_id
GROUP BY c.collection_id, p.product_id, img.image_id, img.image_url
LIMIT $1 OFFSET $2;

-- name: UpdateCollectionWith :one
UPDATE collections
SET 
    name = COALESCE(sqlc.narg('name'), name), 
    image_url = COALESCE(sqlc.narg('image_url'), image_url), 
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = now()
WHERE collection_id = $1
RETURNING *;


-- name: DeleteCollection :exec
DELETE FROM collections WHERE collection_id = $1;

-- name: CountCollections :one
SELECT count(*)
FROM collections
WHERE collection_id = COALESCE(sqlc.narg('collection_id'), collection_id);

-- name: SeedCollections :copyfrom
INSERT INTO collections (name, description, image_url) VALUES ($1, $2, $3);
