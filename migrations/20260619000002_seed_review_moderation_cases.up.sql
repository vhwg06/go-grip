INSERT INTO reviews (id, product_id, order_id, user_id, username, rating, comment, status, created_at, updated_at)
VALUES
    (
        930002,
        'b1111111-1111-1111-1111-111111111111',
        'test-order-0001',
        '22222222-2222-2222-2222-222222222222',
        'test_buyer',
        4,
        'Pending seeded review A',
        'PENDING',
        NOW(),
        NOW()
    ),
    (
        930003,
        'b1111111-1111-1111-1111-111111111111',
        'test-order-0001',
        '22222222-2222-2222-2222-222222222222',
        'test_buyer',
        3,
        'Pending seeded review B',
        'PENDING',
        NOW(),
        NOW()
    )
ON CONFLICT DO NOTHING;
