CREATE TABLE IF NOT EXISTS media_assets (
    id UUID PRIMARY KEY,
    file_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    url TEXT NOT NULL,
    alt_text TEXT,
    owner_type TEXT,
    owner_id UUID,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT media_assets_mime_type_check CHECK (mime_type IN ('image/jpeg', 'image/png', 'image/webp')),
    CONSTRAINT media_assets_size_check CHECK (size_bytes <= 5242880)
);
