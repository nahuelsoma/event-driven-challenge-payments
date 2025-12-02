-- Rollback: Drop Payment Tables

DROP INDEX IF EXISTS idx_payments_idempotency_key;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_user_id;
DROP TABLE IF EXISTS payments;

DROP INDEX IF EXISTS idx_payment_events_payment_id;
DROP TABLE IF EXISTS payment_events;

