INSERT INTO operation_type (name) VALUES 
    ('DEPOSIT'),
    ('WITHDRAW');
	
INSERT INTO users (id, name) VALUES 
    (gen_random_uuid(), 'Alice'),
    (gen_random_uuid(), 'Bob'),
    (gen_random_uuid(), 'Charlie');

INSERT INTO wallets (id, user_id, balance)
SELECT gen_random_uuid(), id, 1000 FROM users;