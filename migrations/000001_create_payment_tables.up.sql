-- Migration: Create Payment Tables (Event Store + Read Model)
-- CQRS Pattern: Event Store is the source of truth, Read Model is for queries

-- EVENT STORE (source of truth, INSERT only, immutable)
CREATE TABLE IF NOT EXISTS payment_events (
    id              TEXT PRIMARY KEY,
    payment_id      TEXT NOT NULL,
    sequence        INTEGER NOT NULL,
    event_type      TEXT NOT NULL,  -- created, reserved, completed, failed
    payload         JSONB NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(payment_id, sequence)  -- Prevents concurrent duplicate events
);

CREATE INDEX IF NOT EXISTS idx_payment_events_payment_id ON payment_events(payment_id);

-- READ MODEL (optimized for queries)
CREATE TABLE IF NOT EXISTS payments (
    id              TEXT PRIMARY KEY,
    idempotency_key TEXT UNIQUE,      -- For request idempotency
    user_id         TEXT NOT NULL,
    amount          DECIMAL(15,2) NOT NULL,
    currency        TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending',
    gateway_ref     TEXT,
    failure_reason  TEXT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Status values: pending, reserved, completed, failed
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_idempotency_key ON payments(idempotency_key);

