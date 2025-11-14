-- name: InsertProductImage :one
INSERT INTO product_images (product_id, image_url, image_id, alt_text, caption, mime_type, file_size, width, height, display_order, is_primary) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: InsertBulkProductImages :copyfrom
INSERT INTO product_images (product_id, image_url, image_id, alt_text, caption, mime_type, file_size, width, height, display_order, is_primary) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetImagesByProductID :many
SELECT * FROM product_images WHERE product_id = $1 ORDER BY display_order;

-- name: GetImageByID :one
SELECT * FROM product_images WHERE id = $1 LIMIT 1;

-- name: GetImageByImageID :one
SELECT * FROM product_images WHERE image_id = $1 LIMIT 1;

-- name: GetPrimaryImageByProductID :one
SELECT * FROM product_images WHERE product_id = $1 AND is_primary = true LIMIT 1;

-- name: GetProductImages :many
SELECT 
    pi.id,
    pi.image_id,
    pi.image_url,
    pi.alt_text,
    pi.caption,
    pi.mime_type,
    pi.file_size,
    pi.width,
    pi.height,
    pi.product_id,
    pi.display_order,
    pi.is_primary
FROM product_images pi WHERE product_id = ANY(sqlc.arg(product_ids)::UUID[])  ORDER BY product_id, display_order;

-- name: UpdateProductImage :exec
UPDATE product_images 
SET 
    image_url = COALESCE(sqlc.narg(image_url), image_url),
    image_id = COALESCE(sqlc.narg(image_id), image_id),
    alt_text = COALESCE(sqlc.narg(alt_text), alt_text),
    caption = COALESCE(sqlc.narg(caption), caption),
    display_order = COALESCE(sqlc.narg(display_order), display_order),
    is_primary = COALESCE(sqlc.narg(is_primary), is_primary),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteProductImage :exec
DELETE FROM product_images WHERE id = $1;

-- name: DeleteProductImagesByProductID :exec
DELETE FROM product_images WHERE product_id = $1;

-- name: SetPrimaryImage :exec
UPDATE product_images SET is_primary = CASE WHEN id = $2 THEN true ELSE false END, updated_at = NOW() WHERE product_id = $1;