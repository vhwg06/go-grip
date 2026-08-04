-- Catalog Base vertical slice storage.
CREATE TABLE IF NOT EXISTS catalog_base_categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    parent_id UUID REFERENCES catalog_base_categories(id),
    position INTEGER NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog_base_attribute_definitions (
    id UUID PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    ordering INTEGER NOT NULL DEFAULT 0,
    value_kind TEXT NOT NULL,
    data_type TEXT NOT NULL DEFAULT '',
    reference_target TEXT NOT NULL DEFAULT '',
    unit_family TEXT NOT NULL DEFAULT '',
    unit TEXT NOT NULL DEFAULT '',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    enum_values JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog_base_masters (
    id UUID PRIMARY KEY,
    kind TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    swatch_media JSONB NOT NULL DEFAULT '[]'::jsonb,
    selling_unit TEXT NOT NULL DEFAULT '',
    quantity NUMERIC,
    base_unit TEXT NOT NULL DEFAULT '',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (kind, name)
);

CREATE TABLE IF NOT EXISTS catalog_base_product_models (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES catalog_base_categories(id),
    description TEXT NOT NULL DEFAULT '',
    warranty_summary JSONB,
    fixed_attributes JSONB NOT NULL DEFAULT '{}'::jsonb,
    fixed_pack_id UUID REFERENCES catalog_base_masters(id),
    measurements JSONB NOT NULL DEFAULT '{}'::jsonb,
    status TEXT NOT NULL DEFAULT 'Draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog_base_product_images (
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES catalog_base_product_models(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    ordering INTEGER NOT NULL DEFAULT 0,
    primary_image BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS catalog_base_variant_dimensions (
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES catalog_base_product_models(id) ON DELETE CASCADE,
    definition_id UUID NOT NULL REFERENCES catalog_base_attribute_definitions(id),
    allowed_values JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (model_id, definition_id)
);

CREATE TABLE IF NOT EXISTS catalog_base_variants (
    id UUID PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES catalog_base_product_models(id) ON DELETE CASCADE,
    selected_options JSONB NOT NULL DEFAULT '{}'::jsonb,
    technical_values JSONB NOT NULL DEFAULT '{}'::jsonb,
    sku TEXT NOT NULL DEFAULT '',
    selling_amount BIGINT,
    selling_currency TEXT NOT NULL DEFAULT '',
    pack_id UUID REFERENCES catalog_base_masters(id),
    status TEXT NOT NULL DEFAULT 'Active',
    canonical_combination TEXT NOT NULL,
    history JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (model_id, canonical_combination)
);

CREATE INDEX IF NOT EXISTS idx_catalog_base_models_status ON catalog_base_product_models(status);
CREATE INDEX IF NOT EXISTS idx_catalog_base_variants_model ON catalog_base_variants(model_id);
CREATE INDEX IF NOT EXISTS idx_catalog_base_variants_sku ON catalog_base_variants(sku);
CREATE INDEX IF NOT EXISTS idx_catalog_base_dimensions_model ON catalog_base_variant_dimensions(model_id);
