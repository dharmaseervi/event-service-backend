package models

import (
	"time"
)

type VendorDeal struct {
	ID              int       `json:"id"`
	VendorID        int       `json:"vendor_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DiscountPercent int       `json:"discount_percent"`
	OriginalPrice   int       `json:"original_price"`
	DealPrice       int       `json:"deal_price"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Photos          []string  `json:"photos"`
	CreatedAt       time.Time `json:"created_at"`
}
