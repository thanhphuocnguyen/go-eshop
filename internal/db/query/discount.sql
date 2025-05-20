-- name: InsertDiscount :one
INSERT INTO discounts 
    (code, description, discount_type, discount_value, min_purchase_amount, max_discount_amount, is_active, usage_limit, used_count, starts_at, expires_at)
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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

-- name: GetDiscount :one
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
WHERE code ILIKE COALESCE(sqlc.narg('code'), discounts.code)
  AND "description" ILIKE COALESCE(sqlc.narg('description'), discounts."description")
  AND discount_type = COALESCE(sqlc.narg('discount_type'), discounts.discount_type)
  AND discount_value = COALESCE(sqlc.narg('discount_value'), discounts.discount_value)
  AND min_purchase_amount = COALESCE(sqlc.narg('min_purchase_amount'), discounts.min_purchase_amount)
  AND max_discount_amount = COALESCE(sqlc.narg('max_discount_amount'), discounts.max_discount_amount)
  AND usage_limit = COALESCE(sqlc.narg('usage_limit'), discounts.usage_limit)
  AND used_count = COALESCE(sqlc.narg('used_count'), discounts.used_count)
  AND is_active = COALESCE(sqlc.narg('is_active'), discounts.is_active)
  AND starts_at >= COALESCE(sqlc.narg('form_date'), discounts.starts_at)
  AND starts_at <= COALESCE(sqlc.narg('to_date'), discounts.starts_at)
ORDER BY id
LIMIT $1
OFFSET $2;


-- name: DeleteDiscount :exec
UPDATE discounts
SET deleted_at = NOW()
WHERE id = $1;

-- name: InsertDiscountProduct :one
INSERT INTO discount_products (discount_id, product_id)
VALUES ($1, $2)
RETURNING id;

-- name: InsertDiscountCategory :one
INSERT INTO discount_categories (discount_id, category_id)
VALUES ($1, $2)
RETURNING id;

-- name: InsertDiscountUser :one
INSERT INTO discount_users (discount_id, user_id)
VALUES ($1, $2)
RETURNING id;

-- name: InsertOrderDiscount :one
INSERT INTO order_discounts (order_id, discount_id, discount_amount)
VALUES ($1, $2, $3)
RETURNING id;
