-- Drop indexes for `cart_items` table
DROP INDEX IF EXISTS idx_cart_items_variant_id;

DROP INDEX IF EXISTS idx_cart_items_cart_id;

-- Drop indexes for `carts` table
DROP INDEX IF EXISTS idx_carts_updated_at;

DROP INDEX IF EXISTS idx_carts_session_id;

DROP INDEX IF EXISTS idx_carts_user_id;

-- Drop tables
DROP TABLE IF EXISTS cart_items;

DROP TABLE IF EXISTS carts;