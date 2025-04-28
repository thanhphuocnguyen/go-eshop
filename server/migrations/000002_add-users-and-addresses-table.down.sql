DROP INDEX IF EXISTS "idx_user_addresses_user_id";

DROP INDEX IF EXISTS "user_addresses_user_id_default_idx";

DROP INDEX IF EXISTS "sessions_user_id_idx";

DROP TABLE IF EXISTS "user_addresses";

DROP TABLE IF EXISTS verify_emails;

DROP TABLE IF EXISTS "sessions" CASCADE;

DROP TABLE IF EXISTS "users" CASCADE;

DROP INDEX IF EXISTS "user_payment_infos_user_id_idx";

DROP TABLE IF EXISTS user_payment_infos;