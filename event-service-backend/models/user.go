package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`        // Local DB primary key
	ClerkID   string    `json:"clerk_id"`  // Unique ID from Clerk
	FullName  string    `json:"full_name"` // Optional if not always needed
	Email     string    `json:"email"`     // Optional if you fetch from Clerk
	Role      string    `json:"role"`      // 'user', 'vendor', 'admin'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
