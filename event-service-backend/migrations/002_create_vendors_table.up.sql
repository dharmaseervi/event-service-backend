CREATE TABLE IF NOT EXISTS vendors (
  id SERIAL PRIMARY KEY,                          -- Auto-incrementing primary key
  vendor_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,  
  title TEXT NOT NULL,                            -- Vendor listing title (e.g., "Royal Palace Hall")
  description TEXT,                               -- Optional longer description
  category TEXT CHECK (
    category IN ('venue', 'catering', 'decor', 'photography')
  ) NOT NULL,                                     -- Category restriction with validation
  price_range TEXT,                               -- e.g., "₹50,000 - ₹2,00,000"
  location TEXT,                                  -- City or area
  photos TEXT[],                                  -- Array of photo URLs (PostgreSQL supports arrays)
  rating NUMERIC(2,1) DEFAULT 0,            -- Rating out of 5, default is 0
  featured BOOLEAN DEFAULT FALSE
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Auto-timestamp
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  -- Auto-timestamp
);
