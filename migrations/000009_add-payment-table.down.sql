DROP INDEX IF EXISTS "payments_order_id_idx";

DROP INDEX IF EXISTS "payments_status_idx";

DROP INDEX IF EXISTS "payments_payment_method_idx";

DROP INDEX IF EXISTS "payments_gateway_payment_intent_id_idx";

DROP INDEX IF EXISTS "payments_gateway_charge_id_idx";

DROP INDEX IF EXISTS "payment_transactions_payment_id_idx";

DROP INDEX IF EXISTS "payment_transactions_status_idx";

DROP INDEX IF EXISTS "payment_transactions_gateway_transaction_id_idx";

DROP TABLE IF EXISTS payment_transactions;

DROP TABLE IF EXISTS payments;