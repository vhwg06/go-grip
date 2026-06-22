ALTER TABLE products
DROP CONSTRAINT IF EXISTS fk_products_intro_article;

DROP INDEX IF EXISTS idx_products_intro_article_id;

ALTER TABLE products
DROP COLUMN IF EXISTS intro_article_id;
