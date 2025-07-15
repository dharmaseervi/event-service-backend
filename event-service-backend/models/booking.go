package models

import (
	"time"
)

type Booking struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VendorID  int       `json:"vendor_id"`
	EventDate time.Time `json:"event_date"`
	Status    string    `json:"status"` // e.g., pending, confirmed, cancelled
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
