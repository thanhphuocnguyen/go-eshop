-- Drop foreign keys
ALTER TABLE "payment_infos" DROP CONSTRAINT payment_infos_user_id_fkey;
ALTER TABLE "category_products" DROP CONSTRAINT category_products_product_id_fkey;
ALTER TABLE "category_products" DROP CONSTRAINT category_products_category_id_fkey;
ALTER TABLE "attribute_values" DROP CONSTRAINT attribute_values_attribute_id_fkey;
ALTER TABLE "shippings" DROP CONSTRAINT shippings_order_id_fkey;
ALTER TABLE "order_items" DROP CONSTRAINT order_items_order_id_fkey;
ALTER TABLE "order_items" DROP CONSTRAINT order_items_product_id_fkey;
ALTER TABLE "orders" DROP CONSTRAINT orders_shipping_id_fkey;
ALTER TABLE "orders" DROP CONSTRAINT orders_user_id_fkey;
ALTER TABLE "cart_items" DROP CONSTRAINT cart_items_cart_id_fkey;
ALTER TABLE "cart_items" DROP CONSTRAINT cart_items_product_id_fkey;
ALTER TABLE "carts" DROP CONSTRAINT carts_user_id_fkey;
ALTER TABLE "sessions" DROP CONSTRAINT "sessions_user_id_fkey";

-- Drop indexes
DROP INDEX IF EXISTS products_price_idx;
DROP INDEX IF EXISTS products_archived_idx;
DROP INDEX IF EXISTS cart_items_product_id_cart_id_idx;
DROP INDEX IF EXISTS orders_status_idx;
DROP INDEX IF EXISTS orders_shipping_id_idx;
DROP INDEX IF EXISTS orders_user_id_idx;
DROP INDEX IF EXISTS orders_user_id_status_idx;
DROP INDEX IF EXISTS order_items_product_id_order_id_idx;
DROP INDEX IF EXISTS shippings_order_id_idx;
DROP INDEX IF EXISTS category_products_category_id_product_id_idx;

-- Drop tables
DROP TABLE IF EXISTS "payment_infos";
DROP TABLE IF EXISTS "category_products";
DROP TABLE IF EXISTS "categories";
DROP TABLE IF EXISTS "attributes";
DROP TABLE IF EXISTS "attribute_values";
DROP TABLE IF EXISTS "shippings";
DROP TABLE IF EXISTS "order_items";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "cart_items";
DROP TABLE IF EXISTS "carts";
DROP TABLE IF EXISTS "products";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "sessions";

-- Drop types
DROP TYPE IF EXISTS "card_type";
DROP TYPE IF EXISTS "payment_type";
DROP TYPE IF EXISTS "payment_status";
DROP TYPE IF EXISTS "order_status";
DROP TYPE IF EXISTS "user_role";
