CREATE TABLE IF NOT EXISTS homepage_blocks (
    id UUID PRIMARY KEY,
    block_type TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    position INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS support_channels (
    id UUID PRIMARY KEY,
    channel_type TEXT NOT NULL,
    label TEXT NOT NULL,
    link TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT true
);
