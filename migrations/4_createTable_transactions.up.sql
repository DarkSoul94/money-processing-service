
CREATE TABLE IF NOT EXISTS transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  type SMALLINT NOT NULL,
  from_account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
  to_account_id INTEGER REFERENCES accounts(id) ON DELETE CASCADE,
  amount DECIMAL(10, 4) NOT NULL
);
