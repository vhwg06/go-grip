ALTER TABLE products
ADD COLUMN IF NOT EXISTS intro_article_id UUID;

CREATE INDEX IF NOT EXISTS idx_products_intro_article_id
ON products(intro_article_id);

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'fk_products_intro_article'
	) THEN
		ALTER TABLE products
		ADD CONSTRAINT fk_products_intro_article
		FOREIGN KEY (intro_article_id)
		REFERENCES content_articles(id)
		ON DELETE SET NULL;
	END IF;
END $$;
