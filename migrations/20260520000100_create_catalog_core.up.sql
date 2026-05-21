CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    sku TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    price BIGINT NOT NULL DEFAULT 0,
    compare_price BIGINT,
    status TEXT NOT NULL DEFAULT 'draft',
    brand TEXT,
    attributes JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id UUID REFERENCES categories(id),
    position INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS product_categories (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE TABLE IF NOT EXISTS product_tags (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, tag_id)
);
