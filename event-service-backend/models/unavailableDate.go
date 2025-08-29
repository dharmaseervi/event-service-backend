package models

import "time"

type UnavailableDate struct {
	ID         int       `json:"id"`
	VendorID   int       `json:"vendor_id"`
	BookedFrom time.Time `json:"booked_from"`
	BookedTo   time.Time `json:"booked_to"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
