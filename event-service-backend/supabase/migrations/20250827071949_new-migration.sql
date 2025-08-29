-- ================================
-- Users table (linked to Clerk)
-- ================================
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  clerk_id TEXT UNIQUE NOT NULL,   -- Clerk’s user ID
  full_name TEXT,
  email TEXT UNIQUE,
  role TEXT DEFAULT 'user',        -- 'user', 'vendor', 'admin'
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_clerk_id ON users(clerk_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- ================================
-- Vendors
-- ================================
CREATE TABLE IF NOT EXISTS vendors (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT,
  category TEXT CHECK (category IN ('venue','catering','decor','photography','entertainment','florist','planning')),
  price_range TEXT,
  location TEXT,
  photos TEXT[],
  rating NUMERIC(2,1) DEFAULT 0,
  featured BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vendors_category ON vendors(category);
CREATE INDEX IF NOT EXISTS idx_vendors_location ON vendors(location);

-- ================================
-- Bookings
-- ================================
CREATE TABLE IF NOT EXISTS bookings (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
  event_date DATE NOT NULL,
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending','confirmed','cancelled','completed')),
  total_price NUMERIC(10,2),
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_vendor_id ON bookings(vendor_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);

-- ================================
-- Deals
-- ================================
CREATE TABLE IF NOT EXISTS deals (
  id SERIAL PRIMARY KEY,
  vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  description TEXT,
  discount_percent INT CHECK (discount_percent >= 0 AND discount_percent <= 100),
  original_price NUMERIC(10,2),
  deal_price NUMERIC(10,2),
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  photos TEXT[],
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deals_vendor_id ON deals(vendor_id);
CREATE INDEX IF NOT EXISTS idx_deals_active ON deals(start_date, end_date);

-- ================================
-- Saved Items (Wishlist/Favorites)
-- ================================
CREATE TABLE IF NOT EXISTS saved_items (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  vendor_id INT REFERENCES vendors(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  UNIQUE(user_id, vendor_id)
);

CREATE INDEX IF NOT EXISTS idx_saved_items_user_id ON saved_items(user_id);

-- ================================
-- Push Tokens (Notifications)
-- ================================
CREATE TABLE IF NOT EXISTS push_tokens (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  token TEXT UNIQUE NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS unavailable_dates (
    id BIGSERIAL PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    booked_from TIMESTAMPTZ NOT NULL,
    booked_to TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_unavailable_dates CHECK (booked_to > booked_from)
);

CREATE INDEX IF NOT EXISTS idx_push_tokens_user_id ON push_tokens(user_id);
