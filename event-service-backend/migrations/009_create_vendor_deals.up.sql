CREATE TABLE IF NOT EXISTS vendor_deals (
  id SERIAL PRIMARY KEY,
  vendor_id INTEGER REFERENCES vendors(id),
  title TEXT NOT NULL,
  description TEXT,
  discount_percent INTEGER,
  original_price INTEGER,
  deal_price INTEGER,
  start_date DATE,
  end_date DATE,
  photos TEXT[],
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
