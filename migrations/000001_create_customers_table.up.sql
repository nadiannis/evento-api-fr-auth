CREATE TABLE IF NOT EXISTS customers (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  password_hash BYTEA NOT NULL,
  balance NUMERIC NOT NULL DEFAULT 0
);

ALTER TABLE customers ADD CONSTRAINT customers_balance_check CHECK (balance >= 0);
