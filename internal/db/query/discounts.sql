-- name: InsertDiscount :one
INSERT INTO discounts 
    (code, description, discount_type, discount_value, min_purchase_amount, max_discount_amount, is_active, usage_limit, starts_at, expires_at)
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: UpdateDiscount :one
UPDATE discounts
SET "description" = COALESCE(sqlc.narg('description'), discounts.description),
    discount_type = COALESCE(sqlc.narg('discount_type'), discounts.discount_type),
    discount_value = COALESCE(sqlc.narg('discount_value'), discounts.discount_value),
    min_purchase_amount = COALESCE(sqlc.narg('min_purchase_amount'), discounts.min_purchase_amount),
    max_discount_amount = COALESCE(sqlc.narg('max_discount_amount'), discounts.max_discount_amount),
    usage_limit = COALESCE(sqlc.narg('usage_limit'), discounts.usage_limit),
    used_count = COALESCE(sqlc.narg('used_count'), discounts.used_count),
    is_active = COALESCE(sqlc.narg('is_active'), discounts.is_active),
    starts_at = COALESCE(sqlc.narg('starts_at'), discounts.starts_at),
    expires_at = COALESCE(sqlc.narg('expires_at'), discounts.expires_at)
WHERE id = $1
RETURNING id;

-- name: GetDiscountByID :one
SELECT *
FROM discounts
WHERE id = $1;

-- name: GetDiscountByCode :one
SELECT *
FROM discounts
WHERE code = $1
LIMIT 1;

-- name: GetDiscounts :many
SELECT id, code, "description", discount_type, discount_value, min_purchase_amount, max_discount_amount, usage_limit, used_count, is_active, starts_at, expires_at
FROM discounts
WHERE 
    (code ILIKE COALESCE(sqlc.narg('search'), discounts.code) OR "description" ILIKE COALESCE(sqlc.narg('search'), discounts."description"))
    AND discount_type = COALESCE(sqlc.narg('discount_type'), discounts.discount_type)
    AND discount_value = COALESCE(sqlc.narg('discount_value'), discounts.discount_value)
    AND min_purchase_amount = COALESCE(sqlc.narg('min_purchase_amount'), discounts.min_purchase_amount)
    AND max_discount_amount = COALESCE(sqlc.narg('max_discount_amount'), discounts.max_discount_amount)
    AND usage_limit = COALESCE(sqlc.narg('usage_limit'), discounts.usage_limit)
    AND used_count = COALESCE(sqlc.narg('used_count'), discounts.used_count)
    AND is_active = COALESCE(sqlc.narg('is_active'), discounts.is_active)
    AND starts_at >= COALESCE(sqlc.narg('from_date'), discounts.starts_at)
    AND starts_at <= COALESCE(sqlc.narg('to_date'), discounts.starts_at)
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CountDiscounts :one
SELECT COUNT(*) FROM discounts;

-- name: DeleteDiscount :exec
UPDATE discounts
SET deleted_at = NOW()
WHERE id = $1;

-- name: InsertDiscountProduct :one
INSERT INTO discount_products (discount_id, product_id)
VALUES ($1, $2)
RETURNING id;

-- name: InsertBulkProductDiscounts :copyfrom
INSERT INTO discount_products (discount_id, product_id)
VALUES ($1, $2);

-- name: DeleteProductDiscountsByDiscountID :exec
DELETE FROM discount_products
WHERE discount_id = $1;

-- name: InsertDiscountCategory :one
INSERT INTO discount_categories (discount_id, category_id)
VALUES ($1, $2)
RETURNING id;

-- name: DeleteCategoryDiscountsByDiscountID :exec
DELETE FROM discount_categories
WHERE discount_id = $1;

-- name: InsertBulkCategoryDiscounts :copyfrom
INSERT INTO discount_categories (discount_id, category_id)
VALUES ($1, $2);

-- name: DeleteDiscountCategory :exec
DELETE FROM discount_categories
WHERE discount_id = $1
  AND category_id = $2;

-- name: InsertDiscountUser :one
INSERT INTO discount_users (discount_id, user_id)
VALUES ($1, $2)
RETURNING id;

-- name: InsertBulkUserDiscounts :copyfrom
INSERT INTO discount_users (discount_id, user_id)
VALUES ($1, $2);

-- name: DeleteUserDiscountsByDiscountID :exec
DELETE FROM discount_users
WHERE discount_id = $1;

-- name: InsertOrderDiscount :one
INSERT INTO order_discounts (order_id, discount_id, discount_amount)
VALUES ($1, $2, $3)
RETURNING id;
