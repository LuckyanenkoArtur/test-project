CREATE VIEW full_wallet_info AS
SELECT 
    u.id AS user_id,
    u.name AS user_name,
    w.id AS wallet_id,
    w.balance
FROM users u
JOIN wallets w ON u.id = w.user_id;

CREATE VIEW transaction_details AS
SELECT 
    t.id AS transaction_id,
    w.id AS wallet_id,
    u.id AS user_id,
    u.name AS user_name,
    o.name AS operation_type,
    t.amount,
    t.created_at
FROM transaction_logs t
JOIN wallets w ON t.wallet_id = w.id
JOIN users u ON w.user_id = u.id
JOIN operation_type o ON t.operation_type_id = o.id;

SELECT * FROM full_wallet_info WHERE user_id = '0c0aa2e0-f3c2-4a10-a940-36baebc2ad52';
SELECT * FROM transaction_details;
