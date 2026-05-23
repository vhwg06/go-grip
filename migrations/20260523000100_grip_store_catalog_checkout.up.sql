BEGIN;

-- Categories and products may already exist from legacy ecommerce migrations.
ALTER TABLE categories
    ADD COLUMN IF NOT EXISTS icon TEXT,
    ADD COLUMN IF NOT EXISTS sort_order INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS name TEXT,
    ADD COLUMN IF NOT EXISTS category TEXT,
    ADD COLUMN IF NOT EXISTS image TEXT,
    ADD COLUMN IF NOT EXISTS is_hot BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS is_shared BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS sort_order INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS purchase_limit INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS purchase_warning TEXT,
    ADD COLUMN IF NOT EXISTS visibility_level INTEGER NOT NULL DEFAULT -1,
    ADD COLUMN IF NOT EXISTS stock_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS locked_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sold_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS rating DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS review_count INTEGER NOT NULL DEFAULT 0;

UPDATE products
SET name = COALESCE(name, title)
WHERE name IS NULL;

CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    product_id TEXT NOT NULL,
    card_key TEXT NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    reserved_order_id TEXT NOT NULL DEFAULT '',
    reserved_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    order_id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL,
    product_name TEXT,
    amount BIGINT NOT NULL,
    email TEXT,
    status TEXT NOT NULL,
    trade_no TEXT,
    card_key TEXT,
    card_ids TEXT,
    paid_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    user_id TEXT,
    username TEXT,
    payee TEXT,
    points_used INTEGER NOT NULL DEFAULT 0,
    quantity INTEGER NOT NULL DEFAULT 1,
    current_payment_id TEXT,
    status_text TEXT,
    status_color TEXT,
    payment_provider_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    provider_payment_id TEXT,
    amount BIGINT NOT NULL,
    status TEXT NOT NULL,
    request_payload_summary TEXT,
    callback_payload_summary TEXT,
    is_signature_valid BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Checkout critical indexes.
CREATE INDEX IF NOT EXISTS idx_cards_reservation_lookup
    ON cards (product_id, is_used, reserved_order_id, reserved_at);
CREATE INDEX IF NOT EXISTS idx_cards_reserved_order_id
    ON cards (reserved_order_id);
CREATE INDEX IF NOT EXISTS idx_orders_status_created_at
    ON orders (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_user_email_created_at
    ON orders (user_id, email, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payments_order_status
    ON payments (order_id, status);
CREATE INDEX IF NOT EXISTS idx_payments_provider_payment_id
    ON payments (provider_payment_id);

COMMIT;
