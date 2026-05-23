BEGIN;

DROP INDEX IF EXISTS idx_cards_product_created_at;
DROP INDEX IF EXISTS idx_refund_requests_status_created_at;
DROP INDEX IF EXISTS idx_admin_messages_target;
DROP INDEX IF EXISTS idx_admin_messages_created_at;

ALTER TABLE refund_requests
    DROP COLUMN IF EXISTS admin_username,
    DROP COLUMN IF EXISTS admin_note,
    DROP COLUMN IF EXISTS processed_at;

DROP TABLE IF EXISTS admin_messages;

COMMIT;
