ALTER TABLE "user_addresses"
DROP CONSTRAINT user_addresses_user_id_fkey;

ALTER TABLE "orders"
DROP CONSTRAINT orders_user_address_id_fkey;

DROP INDEX IF EXISTS "idx_user_addresses_user_id";

DROP TABLE IF EXISTS "user_addresses";