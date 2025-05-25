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
SELECT id, code, "description", discount_type, discount_value,
    min_purchase_amount, max_discount_amount, 
    usage_limit, used_count, is_active, starts_at, expires_at
FROM discounts
WHERE 
    discount_type = COALESCE(sqlc.narg('discount_type'), discounts.discount_type)
    AND is_active = COALESCE(sqlc.narg('is_active'), discounts.is_active)
    AND starts_at >= COALESCE(sqlc.narg('from_date'), discounts.starts_at)
    AND starts_at <= COALESCE(sqlc.narg('to_date'), discounts.starts_at)
    AND code ILIKE '%' || COALESCE(sqlc.narg('search'), discounts.code) || '%'
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetDiscountUsages :many
SELECT usage_limit, used_count, discount_amount, customer_name, order_id, total_price, order_discounts.created_at
FROM discounts
JOIN order_discounts ON discounts.id = order_discounts.discount_id
JOIN orders ON order_discounts.order_id = orders.id
WHERE discounts.id = $1
  AND orders.status IN ('completed', 'confirmed');

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

-- name: GetDiscountProducts :many
SELECT dp.id, dp.discount_id, dp.product_id, p.name, p.base_price, d.discount_type, d.discount_value
FROM discount_products dp
JOIN products p ON dp.product_id = p.id
JOIN discounts d ON dp.discount_id = d.id
WHERE dp.discount_id = $1
ORDER BY dp.id
LIMIT $2
OFFSET $3;

-- name: CountDiscountProducts :one
SELECT COUNT(*) FROM discount_products
WHERE discount_id = $1;

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

-- name: GetDiscountCategories :many
SELECT dc.id, dc.discount_id, dc.category_id, c.name
FROM discount_categories dc
JOIN categories c ON dc.category_id = c.id
WHERE dc.discount_id = $1
ORDER BY dc.id
LIMIT $2
OFFSET $3;

-- name: CountDiscountCategories :one
SELECT COUNT(*) FROM discount_categories
WHERE discount_id = $1;

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

-- name: GetDiscountUsers :many
SELECT du.id, du.discount_id, du.user_id, u.fullname, u.username
FROM discount_users du
JOIN users u ON du.user_id = u.id
WHERE du.discount_id = $1
ORDER BY du.id
LIMIT $2
OFFSET $3;

-- name: CountDiscountUsers :one
SELECT COUNT(*) FROM discount_users
WHERE discount_id = $1;

-- name: DeleteUserDiscountsByDiscountID :exec
DELETE FROM discount_users
WHERE discount_id = $1;

-- name: InsertOrderDiscount :one
INSERT INTO order_discounts (order_id, discount_id, discount_amount)
VALUES ($1, $2, $3)
RETURNING id;
