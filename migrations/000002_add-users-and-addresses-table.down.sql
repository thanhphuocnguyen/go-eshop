

-- Drop indexes first
DROP INDEX IF EXISTS "user_addresses_user_id_default_idx";
DROP INDEX IF EXISTS "sessions_user_id_idx";
DROP INDEX IF EXISTS "idx_user_addresses_user_id";
DROP INDEX IF EXISTS "user_payment_infos_user_id_idx";

-- Drop tables in proper dependency order
-- Tables that reference users must be dropped first
DROP TABLE IF EXISTS "user_addresses";
DROP TABLE IF EXISTS email_verifications;
DROP TABLE IF EXISTS "user_sessions" CASCADE;
DROP TABLE IF EXISTS user_payment_infos;

-- Drop users table last (parent table)
DROP TABLE IF EXISTS "users" CASCADE;