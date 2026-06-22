ALTER TABLE content_articles DROP COLUMN IF EXISTS image_url;
ALTER TABLE content_articles DROP COLUMN IF EXISTS tags;
ALTER TABLE content_articles DROP COLUMN IF EXISTS topic;
ALTER TABLE content_articles DROP COLUMN IF EXISTS priority;

ALTER TABLE static_pages DROP COLUMN IF EXISTS gallery;
