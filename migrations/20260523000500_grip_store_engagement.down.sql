BEGIN;

DROP INDEX IF EXISTS idx_broadcast_messages_created;
DROP INDEX IF EXISTS idx_user_notifications_user_created;
DROP INDEX IF EXISTS idx_reviews_product_id;
DROP INDEX IF EXISTS idx_wishlist_items_created_at;

DROP INDEX IF EXISTS uq_broadcast_reads_message_user;
DROP INDEX IF EXISTS uq_reviews_order_id;
DROP INDEX IF EXISTS uq_wishlist_votes_item_user;

DROP TABLE IF EXISTS broadcast_reads;
DROP TABLE IF EXISTS broadcast_messages;
DROP TABLE IF EXISTS user_notifications;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS wishlist_votes;
DROP TABLE IF EXISTS wishlist_items;

COMMIT;
