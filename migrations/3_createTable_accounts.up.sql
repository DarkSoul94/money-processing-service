
CREATE TABLE IF NOT EXISTS accounts (
  id SERIAL PRIMARY KEY,
  client_id INTEGER NOT NULL REFERENCES clients(id),
  currency_id INTEGER NOT NULL REFERENCES currencys(id),
  ballance DECIMAL(10, 4) NOT NULL
);
