DROP INDEX IF EXISTS "payments_order_id_idx";
DROP INDEX IF EXISTS "order_items_product_id_order_id_idx";
DROP INDEX IF EXISTS "orders_user_id_status_idx";
DROP INDEX IF EXISTS "orders_user_id_idx";
DROP INDEX IF EXISTS "orders_status_idx";

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS "order_items";
DROP TABLE IF EXISTS "orders";
