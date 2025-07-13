package models

import "time"

type SavedVendor struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`   // assumed to be a string (e.g., from Clerk or Firebase)
	VendorID  int       `json:"vendor_id"` // this references `vendors.id`
	CreatedAt time.Time `json:"created_at"`
}
