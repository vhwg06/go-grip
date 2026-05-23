BEGIN;

DROP INDEX IF EXISTS idx_payments_provider_payment_id;
DROP INDEX IF EXISTS idx_payments_order_status;
DROP INDEX IF EXISTS idx_orders_user_email_created_at;
DROP INDEX IF EXISTS idx_orders_status_created_at;
DROP INDEX IF EXISTS idx_cards_reserved_order_id;
DROP INDEX IF EXISTS idx_cards_reservation_lookup;

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS settings;

ALTER TABLE products
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS category,
    DROP COLUMN IF EXISTS image,
    DROP COLUMN IF EXISTS is_hot,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS is_shared,
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS purchase_limit,
    DROP COLUMN IF EXISTS purchase_warning,
    DROP COLUMN IF EXISTS visibility_level,
    DROP COLUMN IF EXISTS stock_count,
    DROP COLUMN IF EXISTS locked_count,
    DROP COLUMN IF EXISTS sold_count,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS review_count;

ALTER TABLE categories
    DROP COLUMN IF EXISTS icon,
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;

COMMIT;
