-- name: CreateImage :one
INSERT INTO images (external_id, url, alt_text, caption, mime_type, file_size, width, height) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: CreateImageAssignment :one
INSERT INTO image_assignments (image_id, entity_id, entity_type, display_order, role) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetProductImagesProductID :many
SELECT * FROM image_assignments
JOIN images ON images.id = image_assignments.image_id
WHERE entity_id = $1 AND entity_type = 'product' ORDER BY display_order;

-- name: GetImageFromID :one
SELECT * FROM images WHERE id = $1 LIMIT 1;

-- name: GetProductImageByExternalID :one
SELECT * FROM images WHERE external_id = $1 LIMIT 1;

-- name: GetProductImageByEntityID :one
SELECT * FROM image_assignments
JOIN images ON images.id = image_assignments.image_id
WHERE entity_id = $1 LIMIT 1;

-- name: UpdateProductImage :exec
UPDATE images 
SET 
    url = COALESCE(sqlc.narg(url), url),
    external_id = COALESCE(sqlc.narg(external_id), external_id) 
WHERE id = $1;

-- name: DeleteProductImage :exec
DELETE FROM images WHERE id = $1;

-- name: DeleteProductImageAssignment :exec
DELETE FROM image_assignments WHERE image_id = $1 AND entity_id = $2 AND entity_type = $3;