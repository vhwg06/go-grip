CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY,
    session_id TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY,
    cart_id UUID NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity >= 1),
    unit_price BIGINT NOT NULL DEFAULT 0,
    product_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    blocked BOOLEAN NOT NULL DEFAULT false
);
