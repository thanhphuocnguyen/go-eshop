-- Drop indexes for `order_items` table
DROP INDEX IF EXISTS idx_order_items_variant_id;

DROP INDEX IF EXISTS idx_order_items_order_id;

-- Drop indexes for `orders` table
DROP INDEX IF EXISTS idx_orders_created_at;

DROP INDEX IF EXISTS idx_orders_customer_email;

DROP INDEX IF EXISTS idx_orders_status;

DROP INDEX IF EXISTS idx_orders_order_date;

DROP INDEX IF EXISTS idx_orders_customer_id;

-- Drop any remaining tables that might reference orders
-- (in case previous migrations didn't run properly)
DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS shipment_items;
DROP TABLE IF EXISTS order_discounts;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_transactions;

-- Update carts table to remove order_id column if it exists
-- (since carts references orders but we don't want to drop carts completely)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'carts' AND column_name = 'order_id'
    ) THEN
        ALTER TABLE carts DROP COLUMN order_id;
    END IF;
END$$;

-- Drop tables (order_items must be dropped before orders due to FK constraint)
DROP TABLE IF EXISTS order_items;

DROP TABLE IF EXISTS orders;
