DROP INDEX IF EXISTS idx_cards_product_created_at;
DROP INDEX IF EXISTS idx_cards_expiry_lookup;
DROP INDEX IF EXISTS idx_cards_product_reserved_at;
DROP INDEX IF EXISTS idx_cards_reserved_order_id;
DROP INDEX IF EXISTS idx_cards_reservation_lookup;
DROP INDEX IF EXISTS uq_daily_checkins_user_day;

DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS daily_checkins_v2;

ALTER TABLE IF EXISTS orders
    DROP COLUMN IF EXISTS points_used;

ALTER TABLE IF EXISTS login_users
    DROP COLUMN IF EXISTS points,
    DROP COLUMN IF EXISTS last_checkin_at,
    DROP COLUMN IF EXISTS consecutive_days;
