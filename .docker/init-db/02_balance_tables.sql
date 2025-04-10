-- Balance Service Tables
CREATE TABLE IF NOT EXISTS account_balances (
    account_id VARCHAR(255) PRIMARY KEY,
    balance DECIMAL(15,2) NOT NULL
);