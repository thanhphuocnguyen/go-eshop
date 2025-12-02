-- name: InsertDiscount :one
INSERT INTO discounts (code, name, description, discount_type, discount_value, min_order_value, max_discount_amount, usage_limit, usage_per_user, times_used, is_active, is_stackable, priority, valid_from, valid_until) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id;

-- name: UpdateDiscount :one
UPDATE discounts SET
    code = COALESCE(sqlc.narg('code'), discounts.code),
    name = COALESCE(sqlc.narg('name'), discounts.name),
    "description" = COALESCE(sqlc.narg('description'), discounts.description),
    discount_type = COALESCE(sqlc.narg('discount_type'), discounts.discount_type),
    discount_value = COALESCE(sqlc.narg('discount_value'), discounts.discount_value),
    min_order_value = COALESCE(sqlc.narg('min_order_value'), discounts.min_order_value),
    max_discount_amount = COALESCE(sqlc.narg('max_discount_amount'), discounts.max_discount_amount),
    usage_limit = COALESCE(sqlc.narg('usage_limit'), discounts.usage_limit),
    is_active = COALESCE(sqlc.narg('is_active'), discounts.is_active),
    is_stackable = COALESCE(sqlc.narg('is_stackable'), discounts.is_stackable),
    priority = COALESCE(sqlc.narg('priority'), discounts.priority),
    valid_from = COALESCE(sqlc.narg('valid_from'), discounts.valid_from),
    valid_until = COALESCE(sqlc.narg('valid_until'), discounts.valid_until),
    usage_per_user = COALESCE(sqlc.narg('usage_per_user'), discounts.usage_per_user),
    updated_at = NOW()
WHERE id = $1 RETURNING *;

-- name: GetDiscountByID :one
SELECT * FROM discounts WHERE id = $1;

-- name: GetDiscountByCodes :many
SELECT discounts.*, JSONB_AGG(DISTINCT discount_rules.*) AS rules
FROM discounts
LEFT JOIN discount_rules ON discounts.id = discount_rules.discount_id
WHERE code = ANY($1)
GROUP BY discounts.id;

-- name: GetDiscountByCode :one
SELECT * FROM discounts WHERE code = $1 LIMIT 1;

-- name: GetDiscounts :many
SELECT * FROM discounts
WHERE 
    discount_type = COALESCE(sqlc.narg('discount_type'), discounts.discount_type)
    AND is_active = COALESCE(sqlc.narg('is_active'), TRUE)
    AND valid_from >= COALESCE(sqlc.narg('from_date'), discounts.valid_from)
    AND valid_from <= COALESCE(sqlc.narg('to_date'), discounts.valid_from)
    AND valid_until >= COALESCE(sqlc.narg('to_date'), discounts.valid_until)
    AND code ILIKE '%' || COALESCE(sqlc.narg('search'), discounts.code) || '%'
LIMIT $1
OFFSET $2;

-- name: GetDiscountUsages :many
SELECT usage_limit, times_used, discount_amount, customer_name, order_id, total_price, discount_usage.created_at
FROM discounts
JOIN discount_usage ON discounts.id = discount_usage.discount_id
JOIN orders ON discount_usage.order_id = orders.id
WHERE discounts.id = $1 AND orders.status IN ('completed', 'confirmed');

-- name: CountDiscounts :one
SELECT COUNT(*) FROM discounts;

-- name: DeleteDiscount :exec
DELETE FROM discounts WHERE id = $1;

-- name: IncrementDiscountUsage :exec
UPDATE discounts SET times_used = times_used + 1 WHERE id = $1;

-- name: DecrementDiscountUsage :exec
UPDATE discounts SET times_used = GREATEST(times_used - 1, 0) WHERE id = $1;


-- name: AddDiscountUsage :one
INSERT INTO discount_usage (discount_id, order_id, user_id, discount_amount) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: RemoveDiscountUsage :exec
DELETE FROM discount_usage WHERE discount_id = $1 AND order_id = $2;

-- name: GetTotalDiscountGiven :one
SELECT SUM(discount_amount) FROM discount_usage WHERE discount_id = $1;

-- name: CountDiscountUsages :one
SELECT COUNT(*) FROM discount_usage WHERE discount_id = $1;

-- name: GetUsersUsingDiscount :many
SELECT DISTINCT order_id FROM discount_usage WHERE discount_id = $1 LIMIT $2 OFFSET $3;

-- name: CountUsersUsingDiscount :one
SELECT COUNT(DISTINCT order_id) FROM discount_usage WHERE discount_id = $1;

-- name: AddDiscountRule :one
INSERT INTO discount_rules (discount_id, rule_type, rule_value) VALUES ($1, $2, $3) RETURNING *;

-- name: GetDiscountRules :many
SELECT * FROM discount_rules WHERE discount_id = $1;

-- name: DeleteDiscountRules :exec
DELETE FROM discount_rules WHERE discount_id = $1;

-- name: GetActiveDiscounts :many
SELECT * FROM discounts WHERE is_active = TRUE AND (valid_from IS NULL OR valid_from <= NOW()) AND (valid_until IS NULL OR valid_until >= NOW());

-- name: GetExpiredDiscounts :many
SELECT * FROM discounts WHERE is_active = TRUE AND valid_until IS NOT NULL AND valid_until < NOW();

-- name: DeactivateDiscount :exec
UPDATE discounts SET is_active = FALSE WHERE id = $1 RETURNING *;

-- name: ReactivateDiscount :exec
UPDATE discounts SET is_active = TRUE WHERE id = $1 RETURNING *;

-- name: GetTopUsedDiscounts :many
SELECT * FROM discounts ORDER BY times_used DESC LIMIT $1 OFFSET $2;
-- name: GetDiscountsByType :many
SELECT * FROM discounts WHERE discount_type = $1 LIMIT $2 OFFSET $3;
-- name: CountDiscountsByType :one
SELECT COUNT(*) FROM discounts WHERE discount_type = $1;
-- name: GetDiscountsByPriority :many
SELECT * FROM discounts WHERE priority = $1 LIMIT $2 OFFSET $3;
-- name: CountDiscountsByPriority :one
SELECT COUNT(*) FROM discounts WHERE priority = $1;

-- name: SeedDiscounts :copyfrom
INSERT INTO discounts (
  code, name, description, discount_type, discount_value,
  min_order_value, max_discount_amount, usage_limit, usage_per_user,
  is_active, is_stackable, priority, valid_from, valid_until) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: InsertDiscountRule :one
INSERT INTO discount_rules (discount_id, rule_type, rule_value) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateDiscountRule :one
UPDATE discount_rules SET
    rule_type = COALESCE(sqlc.narg('rule_type'), discount_rules.rule_type),
    rule_value = COALESCE(sqlc.narg('rule_value'), discount_rules.rule_value)
WHERE id = $1 RETURNING *;

-- name: DeleteDiscountRule :exec
DELETE FROM discount_rules WHERE id = $1;
-- name: GetDiscountRuleByID :one
SELECT * FROM discount_rules WHERE id = $1;
-- name: CountDiscountRules :one
SELECT COUNT(*) FROM discount_rules WHERE discount_id = $1;
-- name: GetDiscountsWithRules :many
SELECT d.*, dr.rule_type, dr.rule_value
FROM discounts d
LEFT JOIN discount_rules dr ON d.id = dr.discount_id
WHERE d.id = COALESCE(sqlc.narg('discount_id'), d.id)
AND dr.rule_type = COALESCE(sqlc.narg('rule_type'), dr.rule_type)
LIMIT $1 OFFSET $2;
-- name: GetActiveDiscountRules :many
SELECT * FROM discount_rules
WHERE discount_id = $1
AND rule_type = COALESCE(sqlc.narg('rule_type'), rule_type)
LIMIT $2 OFFSET $3;
-- name: DeleteDiscountRulesByType :exec
DELETE FROM discount_rules WHERE discount_id = $1 AND rule_type = $2;
-- name: CountDiscountRulesByType :one
SELECT COUNT(*) FROM discount_rules WHERE discount_id = $1 AND rule_type = $2;

-- name: GetAvailableDiscountsForUser :many
SELECT d.* 
FROM discounts d
LEFT JOIN (
  SELECT 
    du.discount_id, 
    COUNT(*) as user_usage_count
  FROM discount_usage du
  JOIN orders o ON du.order_id = o.id
  WHERE o.user_id = $1
  GROUP BY du.discount_id
) user_usage ON d.id = user_usage.discount_id
WHERE d.is_active = TRUE
  AND (d.valid_from IS NULL OR d.valid_from <= NOW())
  AND (d.valid_until IS NULL OR d.valid_until >= NOW())
  AND (d.usage_limit IS NULL OR d.times_used < d.usage_limit)
  AND (d.usage_per_user IS NULL OR COALESCE(user_usage.user_usage_count, 0) < d.usage_per_user)
ORDER BY d.priority DESC, d.discount_value DESC
LIMIT $2 OFFSET $3;

-- name: CountDiscountUsageByDiscountAndUser :one
SELECT COUNT(*) FROM discount_usage du
JOIN orders o ON du.order_id = o.id
WHERE du.discount_id = $1 AND o.user_id = $2;

-- name: CountAvailableDiscountsForUser :one
SELECT COUNT(*)
FROM discounts d
WHERE d.is_active = TRUE
  AND (d.valid_from IS NULL OR d.valid_from <= NOW())
  AND (d.valid_until IS NULL OR d.valid_until >= NOW())
  AND (d.usage_limit IS NULL OR d.times_used < d.usage_limit)
  AND d.id NOT IN (
    SELECT DISTINCT du.discount_id 
    FROM discount_usage du 
    JOIN orders o ON du.order_id = o.id 
    WHERE o.user_id = $1
  );