BEGIN;

CREATE TABLE IF NOT EXISTS wishlist_items (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    user_id TEXT NOT NULL,
    username TEXT,
    vote_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wishlist_votes (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    product_id TEXT NOT NULL,
    order_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT,
    rating INTEGER NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    title_key TEXT NOT NULL,
    content_key TEXT NOT NULL,
    data TEXT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS broadcast_messages (
    id BIGSERIAL PRIMARY KEY,
    title_key TEXT NOT NULL,
    content_key TEXT NOT NULL,
    data TEXT,
    sender TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS broadcast_reads (
    id BIGSERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL,
    user_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_wishlist_votes_item_user
    ON wishlist_votes (item_id, user_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_reviews_order_id
    ON reviews (order_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_broadcast_reads_message_user
    ON broadcast_reads (message_id, user_id);

CREATE INDEX IF NOT EXISTS idx_wishlist_items_created_at
    ON wishlist_items (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reviews_product_id
    ON reviews (product_id);
CREATE INDEX IF NOT EXISTS idx_user_notifications_user_created
    ON user_notifications (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_broadcast_messages_created
    ON broadcast_messages (created_at DESC);

COMMIT;
