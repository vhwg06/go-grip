BEGIN;

DROP INDEX IF EXISTS idx_cards_expiry_lookup;
DROP INDEX IF EXISTS idx_cards_product_reserved_at;
DROP INDEX IF EXISTS idx_orders_status_updated_at;
DROP INDEX IF EXISTS idx_orders_pending_created_at;
DROP INDEX IF EXISTS uq_payments_provider_payment_id;
DROP INDEX IF EXISTS uq_payments_idempotency_key;

ALTER TABLE cards
    DROP COLUMN IF EXISTS reserved_at;

ALTER TABLE orders
    DROP COLUMN IF EXISTS failed_at;

ALTER TABLE payments
    DROP COLUMN IF EXISTS idempotency_key;

COMMIT;
