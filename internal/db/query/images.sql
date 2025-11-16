-- name: InsertProductImage :one
INSERT INTO product_images (product_id, image_url, image_id, alt_text, caption, display_order) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: InsertBulkProductImages :copyfrom
INSERT INTO product_images (product_id, image_url, image_id, alt_text, caption, display_order) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetImagesByProductID :many
SELECT * FROM product_images WHERE product_id = $1 ORDER BY display_order;

-- name: GetImageByID :one
SELECT * FROM product_images WHERE id = $1 LIMIT 1;

-- name: GetImageByImageID :one
SELECT * FROM product_images WHERE image_id = $1 LIMIT 1;

-- name: GetPrimaryImageByProductID :one
SELECT * FROM product_images WHERE product_id = $1 LIMIT 1;

-- name: GetProductImages :many
SELECT  p_img.id, p_img.image_id, p_img.image_url, p_img.alt_text, p_img.caption, p_img.product_id, p_img.display_order FROM product_images p_img WHERE product_id = ANY(sqlc.arg(product_ids)::UUID[])  ORDER BY product_id, display_order;
-- name: UpdateProductImage :exec
UPDATE product_images 
SET 
    image_url = COALESCE(sqlc.narg(image_url), image_url),
    image_id = COALESCE(sqlc.narg(image_id), image_id),
    alt_text = COALESCE(sqlc.narg(alt_text), alt_text),
    caption = COALESCE(sqlc.narg(caption), caption),
    display_order = COALESCE(sqlc.narg(display_order), display_order),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteProductImage :exec
DELETE FROM product_images WHERE id = $1;

-- name: DeleteProductImagesByProductID :exec
DELETE FROM product_images WHERE product_id = $1;
