BEGIN;

CREATE TABLE IF NOT EXISTS admin_messages (
    id BIGSERIAL PRIMARY KEY,
    target_type TEXT NOT NULL,
    target_value TEXT,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    sender TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE refund_requests
    ADD COLUMN IF NOT EXISTS admin_username TEXT,
    ADD COLUMN IF NOT EXISTS admin_note TEXT,
    ADD COLUMN IF NOT EXISTS processed_at TIMESTAMPTZ;

ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_admin_messages_created_at
    ON admin_messages (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_messages_target
    ON admin_messages (target_type, target_value);
CREATE INDEX IF NOT EXISTS idx_refund_requests_status_created_at
    ON refund_requests (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_cards_product_created_at
    ON cards (product_id, created_at DESC);

COMMIT;
