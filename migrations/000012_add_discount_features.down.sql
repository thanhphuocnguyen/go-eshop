-- Remove discounted_price column from order_items table
ALTER TABLE order_items DROP COLUMN IF EXISTS discounted_price;

-- Drop indexes
DROP INDEX IF EXISTS idx_order_discounts_discount_id;
DROP INDEX IF EXISTS idx_order_discounts_order_id;
DROP INDEX IF EXISTS idx_discount_users_user_id;
DROP INDEX IF EXISTS idx_discount_users_discount_id;
DROP INDEX IF EXISTS idx_discount_categories_category_id;
DROP INDEX IF EXISTS idx_discount_categories_discount_id;
DROP INDEX IF EXISTS idx_discount_products_product_id;
DROP INDEX IF EXISTS idx_discount_products_discount_id;
DROP INDEX IF EXISTS idx_discounts_dates;
DROP INDEX IF EXISTS idx_discounts_is_active;
DROP INDEX IF EXISTS idx_discounts_code;

-- Drop trigger and function


-- Drop tables
DROP TABLE IF EXISTS order_discounts;
DROP TABLE IF EXISTS discount_users;
DROP TABLE IF EXISTS discount_categories;
DROP TABLE IF EXISTS discount_products;
DROP TABLE IF EXISTS discounts;

-- Drop trigger and function (after dropping discounts table)
DROP TRIGGER IF EXISTS update_discount_timestamp ON discounts;
DROP FUNCTION IF EXISTS update_discount_updated_at();
