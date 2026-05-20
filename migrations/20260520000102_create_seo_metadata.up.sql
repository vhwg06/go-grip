CREATE TABLE IF NOT EXISTS seo_metadata (
    id UUID PRIMARY KEY,
    owner_type TEXT NOT NULL,
    owner_id UUID NOT NULL,
    meta_title TEXT,
    meta_description TEXT,
    slug TEXT NOT NULL,
    alt_text TEXT,
    UNIQUE (owner_type, owner_id),
    UNIQUE (owner_type, slug)
);
