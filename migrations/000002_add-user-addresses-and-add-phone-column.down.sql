ALTER TABLE "user_addresses"
DROP CONSTRAINT user_addresses_user_id_fkey;

ALTER TABLE "users"
DROP COLUMN "phone";

DROP TABLE IF EXISTS "user_addresses";