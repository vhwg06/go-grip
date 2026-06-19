ALTER TABLE reviews
    ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'PENDING';

CREATE INDEX IF NOT EXISTS idx_reviews_status_created_at
    ON reviews (status, created_at DESC);

UPDATE reviews
SET status = 'APPROVED'
WHERE status IS NULL OR status = '';
