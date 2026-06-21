BEGIN;

-- Core legacy users for task/history modules.
INSERT INTO users (id, username, email, password_hash, display_name, role_id, status, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'test_admin', 'test_admin@example.com', '$2a$10$aGaIFRHg60wJ36Gxvcal1OS/7UNyBLTZrPPGAeeGDuQBtjqWw829O', 'Test Admin', '00000000-0000-0000-0000-000000000001', 'active', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'test_buyer', 'test_buyer@example.com', '$2a$10$aGaIFRHg60wJ36Gxvcal1OS/7UNyBLTZrPPGAeeGDuQBtjqWw829O', 'Test Buyer', '00000000-0000-0000-0000-000000000005', 'active', NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO login_users (
    id, provider, provider_id, username, email, password_hash, role_id, role, status,
    points, trust_level, is_admin, desktop_notifications_enabled, created_at, updated_at
)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'local', 'test_admin_local', 'test_admin', 'test_admin@example.com', '$2a$10$aGaIFRHg60wJ36Gxvcal1OS/7UNyBLTZrPPGAeeGDuQBtjqWw829O', '00000000-0000-0000-0000-000000000001', 'Administrator', 'active', 5000, 10, TRUE, TRUE, NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'local', 'test_buyer_local', 'test_buyer', 'test_buyer@example.com', '$2a$10$aGaIFRHg60wJ36Gxvcal1OS/7UNyBLTZrPPGAeeGDuQBtjqWw829O', '00000000-0000-0000-0000-000000000005', 'Subscriber', 'active', 1200, 0, FALSE, FALSE, NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO tasks (id, user_id, title, description, status, created_at, updated_at)
VALUES
    ('aaaa0000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'Test task', 'Smoke test task item', 'todo', NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO history (id, user_id, source, destination, original, translation)
VALUES
    (9001, '22222222-2222-2222-2222-222222222222', 'en', 'vi', 'hello', 'xin chao')
ON CONFLICT DO NOTHING;

-- Catalog.
INSERT INTO categories (id, name, parent_id, position, is_active, icon, sort_order, created_at, updated_at)
VALUES
    ('a1111111-1111-1111-1111-111111111111', 'Test Gift Cards', NULL, 1, TRUE, 'gift', 1, NOW(), NOW()),
    ('a2222222-2222-2222-2222-222222222222', 'Test Topups', NULL, 2, TRUE, 'bolt', 2, NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO tags (id, name)
VALUES
    ('c1111111-1111-1111-1111-111111111111', 'test-hot'),
    ('c2222222-2222-2222-2222-222222222222', 'test-sale')
ON CONFLICT DO NOTHING;

INSERT INTO products (
    id, title, sku, description, price, compare_price, status, brand, attributes, created_at, updated_at,
    name, category, image, is_hot, is_active, is_shared, sort_order, purchase_limit, purchase_warning,
    visibility_level, stock_count, locked_count, sold_count, rating, review_count
)
VALUES
    (
        'b1111111-1111-1111-1111-111111111111', 'Test Shared Card 100K', 'TEST-SHARED-100K',
        'Shared card for API smoke tests', 100000, 120000, 'active', 'TEST_BRAND', '{"region":"VN"}'::jsonb, NOW(), NOW(),
        'Test Shared Card 100K', 'a1111111-1111-1111-1111-111111111111', 'https://example.com/test-shared.png',
        TRUE, TRUE, TRUE, 10, 2, 'Only 2 per account', -1, 50, 0, 10, 4.8, 12
    ),
    (
        'b2222222-2222-2222-2222-222222222222', 'Test Unique Card 50K', 'TEST-UNIQUE-50K',
        'Unique card for checkout tests', 50000, 65000, 'active', 'TEST_BRAND', '{"region":"VN"}'::jsonb, NOW(), NOW(),
        'Test Unique Card 50K', 'a1111111-1111-1111-1111-111111111111', 'https://example.com/test-unique.png',
        FALSE, TRUE, FALSE, 20, 1, 'Single purchase only', -1, 20, 0, 5, 4.5, 4
    )
ON CONFLICT DO NOTHING;

INSERT INTO product_categories (product_id, category_id)
VALUES
    ('b1111111-1111-1111-1111-111111111111', 'a1111111-1111-1111-1111-111111111111'),
    ('b2222222-2222-2222-2222-222222222222', 'a1111111-1111-1111-1111-111111111111')
ON CONFLICT (product_id, category_id) DO NOTHING;

INSERT INTO product_tags (product_id, tag_id)
VALUES
    ('b1111111-1111-1111-1111-111111111111', 'c1111111-1111-1111-1111-111111111111'),
    ('b2222222-2222-2222-2222-222222222222', 'c2222222-2222-2222-2222-222222222222')
ON CONFLICT (product_id, tag_id) DO NOTHING;

INSERT INTO cards (id, product_id, card_key, is_used, reserved_order_id, created_at)
VALUES
    (900001, 'b2222222-2222-2222-2222-222222222222', 'TEST-CARD-0001', FALSE, '', NOW()),
    (900002, 'b2222222-2222-2222-2222-222222222222', 'TEST-CARD-0002', FALSE, '', NOW()),
    (900003, 'b2222222-2222-2222-2222-222222222222', 'TEST-CARD-0003', FALSE, '', NOW())
ON CONFLICT DO NOTHING;

INSERT INTO settings (key, value, updated_at)
VALUES
    ('test.announcement', 'Test announcement for migration seed', NOW()),
    ('test.support.email', 'support-test@example.com', NOW())
ON CONFLICT (key) DO NOTHING;

-- Checkout and order lifecycle.
INSERT INTO orders (
    order_id, product_id, product_name, amount, email, status, trade_no, card_key, card_ids, paid_at, delivered_at,
    user_id, username, payee, points_used, quantity, current_payment_id, status_text, status_color, payment_provider_id,
    created_at, updated_at
)
VALUES
    (
        'test-order-0001', 'b1111111-1111-1111-1111-111111111111', 'Test Shared Card 100K', 100000,
        'test_buyer@example.com', 'delivered', 'TEST-TRADE-0001', 'TEST-SHARED-LIVE-KEY', '',
        NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', '22222222-2222-2222-2222-222222222222', 'test_buyer',
        'test_payee', 0, 1, 'test-pay-0001', 'Da giao hang', '#10b981', 'provider-payment-0001', NOW() - INTERVAL '2 days', NOW() - INTERVAL '1 day'
    ),
    (
        'test-order-0002', 'b2222222-2222-2222-2222-222222222222', 'Test Unique Card 50K', 50000,
        'test_buyer@example.com', 'cancelled', '', '', '900001',
        NULL, NULL, '22222222-2222-2222-2222-222222222222', 'test_buyer',
        'test_payee', 0, 1, '', 'Da huy', '#ef4444', '', NOW(), NOW()
    )
ON CONFLICT DO NOTHING;

INSERT INTO payments (
    id, order_id, provider, provider_payment_id, amount, status, request_payload_summary,
    callback_payload_summary, is_signature_valid, processed_at, created_at, idempotency_key
)
VALUES
    (
        'test-pay-0001', 'test-order-0001', 'test_gateway', 'provider-payment-0001', 100000, 'paid',
        '{"source":"seed"}', '{"status":"ok"}', TRUE, NOW() - INTERVAL '1 day', NOW() - INTERVAL '2 days', 'test-idempotency-0001'
    )
ON CONFLICT DO NOTHING;

INSERT INTO refund_requests (
    id, order_id, user_id, username, reason, status, admin_username, admin_note, processed_at, created_at, updated_at
)
VALUES
    (910001, 'test-order-0001', '22222222-2222-2222-2222-222222222222', 'test_buyer', 'Need refund test flow', 'pending', '', '', NULL, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Engagement.
INSERT INTO wishlist_items (id, title, description, user_id, username, vote_count, created_at, updated_at)
VALUES
    (920001, 'Test wishlist item', 'Wishlist seeded by migration', '22222222-2222-2222-2222-222222222222', 'test_buyer', 1, NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO wishlist_votes (id, item_id, user_id, created_at)
VALUES
    (920101, 920001, '11111111-1111-1111-1111-111111111111', NOW())
ON CONFLICT DO NOTHING;

INSERT INTO reviews (id, product_id, order_id, user_id, username, rating, comment, created_at, updated_at)
VALUES
    (930001, 'b1111111-1111-1111-1111-111111111111', 'test-order-0001', '22222222-2222-2222-2222-222222222222', 'test_buyer', 5, 'Seeded review for smoke test', NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO user_notifications (id, user_id, type, title_key, content_key, data, is_read, created_at)
VALUES
    (940001, '22222222-2222-2222-2222-222222222222', 'order_status', 'notif.order.paid.title', 'notif.order.paid.body', '{"orderId":"test-order-0001"}', FALSE, NOW())
ON CONFLICT DO NOTHING;

INSERT INTO broadcast_messages (id, title_key, content_key, data, sender, created_at)
VALUES
    (950001, 'broadcast.test.title', 'broadcast.test.body', '{"campaign":"seed"}', 'test_admin', NOW())
ON CONFLICT DO NOTHING;

INSERT INTO broadcast_reads (id, message_id, user_id, created_at)
VALUES
    (950101, 950001, '22222222-2222-2222-2222-222222222222', NOW())
ON CONFLICT DO NOTHING;

INSERT INTO admin_messages (id, target_type, target_value, title, body, sender, created_at)
VALUES
    (960001, 'broadcast', NULL, 'Seed message', 'Admin seeded message', 'test_admin', NOW())
ON CONFLICT DO NOTHING;

INSERT INTO daily_checkins_v2 (id, user_id, checkin_date, reward, streak_after, created_at)
VALUES
    (970001, '22222222-2222-2222-2222-222222222222', NOW()::date, 100, 3, NOW())
ON CONFLICT DO NOTHING;

-- Cart and lead flow fixtures.
INSERT INTO carts (id, session_id, status, created_at, updated_at)
VALUES
    ('d1111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'active', NOW(), NOW())
ON CONFLICT DO NOTHING;

INSERT INTO cart_items (id, cart_id, product_id, quantity, unit_price, product_snapshot, blocked)
VALUES
    ('e1111111-1111-1111-1111-111111111111', 'd1111111-1111-1111-1111-111111111111', 'b2222222-2222-2222-2222-222222222222', 1, 50000, '{"title":"Test Unique Card 50K"}'::jsonb, FALSE)
ON CONFLICT DO NOTHING;

INSERT INTO order_requests (
    id, cart_id, customer_name, customer_phone, customer_email, address, note, status, cart_snapshot, created_at
)
VALUES
    (
        'f1111111-1111-1111-1111-111111111111', 'd1111111-1111-1111-1111-111111111111',
        'Test Buyer', '0900000000', 'test_buyer@example.com', 'HCMC', 'Seeded order request', 'new',
        '{"items":[{"product_id":"b2222222-2222-2222-2222-222222222222","quantity":1}]}'::jsonb, NOW()
    )
ON CONFLICT DO NOTHING;

INSERT INTO lead_submissions (
    id, source, customer_name, customer_phone, customer_email, message, status, created_at
)
VALUES
    (
        'f2222222-2222-2222-2222-222222222222', 'landing_page', 'Seed Lead',
        '0911000000', 'lead@example.com', 'Interested in pricing', 'new', NOW()
    )
ON CONFLICT (id) DO NOTHING;

-- Keep serial sequences above explicit IDs.
SELECT setval(pg_get_serial_sequence('history', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM history), 1), TRUE);
SELECT setval(pg_get_serial_sequence('cards', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM cards), 1), TRUE);
SELECT setval(pg_get_serial_sequence('refund_requests', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM refund_requests), 1), TRUE);
SELECT setval(pg_get_serial_sequence('wishlist_items', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM wishlist_items), 1), TRUE);
SELECT setval(pg_get_serial_sequence('wishlist_votes', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM wishlist_votes), 1), TRUE);
SELECT setval(pg_get_serial_sequence('reviews', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM reviews), 1), TRUE);
SELECT setval(pg_get_serial_sequence('user_notifications', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM user_notifications), 1), TRUE);
SELECT setval(pg_get_serial_sequence('broadcast_messages', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM broadcast_messages), 1), TRUE);
SELECT setval(pg_get_serial_sequence('broadcast_reads', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM broadcast_reads), 1), TRUE);
SELECT setval(pg_get_serial_sequence('admin_messages', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM admin_messages), 1), TRUE);
SELECT setval(pg_get_serial_sequence('daily_checkins_v2', 'id'), GREATEST((SELECT COALESCE(MAX(id), 1) FROM daily_checkins_v2), 1), TRUE);

COMMIT;
