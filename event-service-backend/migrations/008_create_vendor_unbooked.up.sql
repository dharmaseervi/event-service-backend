CREATE TABLE IF NOT EXISTS vendor_bookings (
  id SERIAL PRIMARY KEY,
  vendor_id INTEGER REFERENCES vendors(id),
  booked_from DATE NOT NULL,
  booked_to DATE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
