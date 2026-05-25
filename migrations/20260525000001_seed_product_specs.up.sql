INSERT INTO product_details (product_id, key, value, sort_order)
VALUES
    ('b1111111-1111-1111-1111-111111111111', 'Material', 'Brass', 0),
    ('b1111111-1111-1111-1111-111111111111', 'Finish', 'Gold', 1),
    ('b2222222-2222-2222-2222-222222222222', 'Material', 'Bronze', 0),
    ('b2222222-2222-2222-2222-222222222222', 'Finish', 'Antique', 1)
ON CONFLICT (product_id, key) DO NOTHING;
