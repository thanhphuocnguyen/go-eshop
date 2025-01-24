-- name: CreateImage :one
INSERT INTO images (product_id, variant_id, image_url, external_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetImageByID :one
SELECT * FROM images WHERE image_id = $1 LIMIT 1;

-- name: GetImageByExternalID :one
SELECT * FROM images WHERE external_id = $1 LIMIT 1;

-- name: GetImageByProductID :one
SELECT * FROM images WHERE product_id = $1 AND variant_id = NULL LIMIT 1;

-- name: GetImageByVariantID :one
SELECT * FROM images WHERE variant_id = $1 LIMIT 1;

-- name: UpdateImage :exec
UPDATE images 
SET 
    product_id = COALESCE(sqlc.narg(product_id), product_id),
    variant_id = COALESCE(sqlc.narg(variant_id), variant_id), 
    image_url = COALESCE(sqlc.narg(image_url), image_url), 
    external_id = COALESCE(sqlc.narg(external_id), external_id) 
WHERE image_id = $1;

-- name: DeleteImage :exec
DELETE FROM images WHERE image_id = $1;