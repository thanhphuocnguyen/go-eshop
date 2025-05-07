-- name: CreateImage :one
INSERT INTO images (external_id, url, alt_text, caption, mime_type, file_size, width, height) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: CreateImageAssignment :one
INSERT INTO image_assignments (image_id, entity_id, entity_type, display_order, role) VALUES ($1, $2, $3, $4, $5) RETURNING *;
-- name: CreateBulkImageAssignments :copyfrom
INSERT INTO image_assignments (image_id, entity_id, entity_type, display_order, role) VALUES ($1, $2, $3, $4, $5);

-- name: GetImagesByEntityID :many
SELECT * FROM image_assignments
JOIN images ON images.id = image_assignments.image_id
WHERE entity_id = $1 ORDER BY display_order;

-- name: GetImageFromID :one
SELECT * FROM images
JOIN image_assignments ON images.id = image_assignments.image_id
WHERE images.id = $1 AND entity_type = $2 LIMIT 1;

-- name: GetImageFromExternalID :one
SELECT * FROM images
JOIN image_assignments ON images.id = image_assignments.image_id
WHERE external_id = $1 LIMIT 1;

-- name: GetProductImageByEntityID :one
SELECT * FROM image_assignments
JOIN images ON images.id = image_assignments.image_id
WHERE entity_id = $1 LIMIT 1;

-- name: GetProductImagesAssigned :many
SELECT 
    images.id,
    images.external_id,
    images.url,
    images.alt_text,
    images.caption,
    images.mime_type,
    images.file_size,
    images.width,
    images.height,
    image_assignments.entity_id,
    image_assignments.entity_type,
    image_assignments.display_order,
    image_assignments.role
FROM images
JOIN image_assignments ON images.id = image_assignments.image_id
WHERE entity_id = ANY(sqlc.arg(entity_ids)::UUID[]) ORDER BY entity_id, display_order;

-- name: UpdateProductImage :exec
UPDATE images 
SET 
    url = COALESCE(sqlc.narg(url), url),
    external_id = COALESCE(sqlc.narg(external_id), external_id) 
WHERE id = $1;

-- name: UpdateProductImageAssignment :exec
UPDATE image_assignments
SET 
    display_order = COALESCE(sqlc.narg(display_order), display_order),
    role = COALESCE(sqlc.narg(role), role)
WHERE image_id = $1 AND entity_id = $2 AND entity_type = $3;

-- name: DeleteProductImage :exec
DELETE FROM images WHERE id = $1;

-- name: DeleteProductImageAssignment :exec
DELETE FROM image_assignments WHERE image_id = $1 AND entity_id = $2 AND entity_type = $3;

-- name: DeleteImageAssignments :exec
DELETE FROM image_assignments WHERE image_id = $1 AND entity_type = $2;