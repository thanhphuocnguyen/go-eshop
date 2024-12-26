-- name: CreateImage :one
INSERT INTO images (product_id, variant_id, image_url, cloudinary_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetImageByID :one
SELECT * FROM images WHERE image_id = $1 LIMIT 1;

-- name: GetImagesByProductID :one
SELECT * FROM images WHERE product_id = $1 AND image_id = $2;

-- name: GetImagesByVariantID :many
SELECT * FROM images WHERE variant_id = $1;

-- name: DeleteImage :exec
DELETE FROM images WHERE image_id = $1;