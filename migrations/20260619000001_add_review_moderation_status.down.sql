DROP INDEX IF EXISTS idx_reviews_status_created_at;

ALTER TABLE reviews
    DROP COLUMN IF EXISTS status;
