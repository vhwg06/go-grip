BEGIN;

ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS idempotency_key TEXT;

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS failed_at TIMESTAMPTZ;

ALTER TABLE cards
    ADD COLUMN IF NOT EXISTS reserved_at TIMESTAMPTZ;

CREATE UNIQUE INDEX IF NOT EXISTS uq_payments_idempotency_key
    ON payments (idempotency_key)
    WHERE idempotency_key IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_payments_provider_payment_id
    ON payments (provider_payment_id)
    WHERE provider_payment_id IS NOT NULL AND provider_payment_id <> '';

CREATE INDEX IF NOT EXISTS idx_orders_pending_created_at
    ON orders (status, created_at)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_orders_status_updated_at
    ON orders (status, updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_cards_product_reserved_at
    ON cards (product_id, reserved_at);

CREATE INDEX IF NOT EXISTS idx_cards_expiry_lookup
    ON cards (expires_at, is_used);

COMMIT;
