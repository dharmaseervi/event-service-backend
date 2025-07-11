package models

import (
	"time"

	"github.com/lib/pq"
)

type VendorListing struct {
	ID          int            `json:"id"`
	VendorID    int            `json:"vendor_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Category    string         `json:"category"` // 'venue', 'catering', 'decor', 'photography'
	PriceRange  string         `json:"price_range"`
	Location    string         `json:"location"`
	Photos      pq.StringArray `json:"photos" gorm:"type:text[]"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
