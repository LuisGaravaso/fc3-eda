-- Populate clients
INSERT INTO clients (id, name, email, created_at) VALUES 
('b7295961-c51c-438a-90e2-78e62f18b726', 'Luis Garavaso', 'j.j@email.com', NOW()),
('31f52dea-7856-42dc-9364-f508fa74d5d7', 'Jane Doe', 'jane.j@email.com', NOW());

-- Populate accounts with initial balance of 100
INSERT INTO accounts (id, client_id, balance, created_at) VALUES 
('7ebc23f5-dd1e-4d93-9490-9fce5052a5f5', 'b7295961-c51c-438a-90e2-78e62f18b726', 100.00, NOW()),
('dff2d137-bba6-4138-81b9-3da7567f122b', '31f52dea-7856-42dc-9364-f508fa74d5d7', 100.00, NOW());

-- Populate balance service database with matching account balances
INSERT INTO account_balances (account_id, balance) VALUES 
('7ebc23f5-dd1e-4d93-9490-9fce5052a5f5', 100.00),
('dff2d137-bba6-4138-81b9-3da7567f122b', 100.00);