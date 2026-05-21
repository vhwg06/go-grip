CREATE TABLE IF NOT EXISTS order_requests (
    id UUID PRIMARY KEY,
    cart_id UUID NOT NULL REFERENCES carts(id),
    customer_name TEXT NOT NULL,
    customer_phone TEXT NOT NULL,
    customer_email TEXT,
    address TEXT,
    note TEXT,
    status TEXT NOT NULL DEFAULT 'new',
    cart_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS lead_submissions (
    id UUID PRIMARY KEY,
    source TEXT NOT NULL,
    customer_name TEXT NOT NULL,
    customer_phone TEXT NOT NULL,
    customer_email TEXT,
    message TEXT,
    status TEXT NOT NULL DEFAULT 'new',
    created_at TIMESTAMPTZ NOT NULL
);
