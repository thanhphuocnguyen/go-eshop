-- Drop indexes for `order_items` table
DROP INDEX IF EXISTS idx_order_items_variant_id;

DROP INDEX IF EXISTS idx_order_items_order_id;

-- Drop indexes for `orders` table
DROP INDEX IF EXISTS idx_orders_created_at;

DROP INDEX IF EXISTS idx_orders_customer_email;

DROP INDEX IF EXISTS idx_orders_status;

DROP INDEX IF EXISTS idx_orders_order_date;

DROP INDEX IF EXISTS idx_orders_customer_id;

-- Drop tables
DROP TABLE IF EXISTS order_items;

DROP TABLE IF EXISTS orders;