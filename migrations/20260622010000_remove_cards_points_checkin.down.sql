ALTER TABLE IF EXISTS login_users
    ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_checkin_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS consecutive_days INTEGER NOT NULL DEFAULT 0;

ALTER TABLE IF EXISTS orders
    ADD COLUMN IF NOT EXISTS points_used INTEGER NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS daily_checkins_v2 (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES login_users(id) ON DELETE CASCADE,
    checkin_date DATE NOT NULL,
    reward INTEGER NOT NULL DEFAULT 0,
    streak_after INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_daily_checkins_user_day
    ON daily_checkins_v2 (user_id, checkin_date);

CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    card_key TEXT NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    reserved_order_id TEXT NOT NULL DEFAULT '',
    reserved_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cards_reservation_lookup
    ON cards (product_id, is_used, reserved_order_id, reserved_at);
CREATE INDEX IF NOT EXISTS idx_cards_reserved_order_id
    ON cards (reserved_order_id);
CREATE INDEX IF NOT EXISTS idx_cards_product_reserved_at
    ON cards (product_id, reserved_at);
CREATE INDEX IF NOT EXISTS idx_cards_expiry_lookup
    ON cards (expires_at, is_used);
CREATE INDEX IF NOT EXISTS idx_cards_product_created_at
    ON cards (product_id, created_at DESC);
