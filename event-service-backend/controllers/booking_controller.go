package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/gin-gonic/gin"
)

// POST /bookings
func CreateBooking(c *gin.Context) {
	var booking models.Booking
	log.Printf("Received request to booking: %v", c.Request.Body)

	if err := c.ShouldBindJSON(&booking); err != nil {
		log.Printf("BindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("booking vendorId: %d", booking.VendorID)
	if booking.UserID == 0 || booking.VendorID == 0 || booking.EventDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, vendor_id, and event_date are required"})
		return
	}

	if booking.Status == "" {
		booking.Status = "pending"
	}
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	// Insert booking
	query := `
		INSERT INTO bookings (user_id, vendor_id, event_date, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	err := config.DB.QueryRow(
		query,
		booking.UserID,
		booking.VendorID,
		booking.EventDate,
		booking.Status,
		booking.Notes,
		booking.CreatedAt,
		booking.UpdatedAt,
	).Scan(&booking.ID)

	if err != nil {
		log.Printf("❌ Insert failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created",
		"booking": booking,
	})
}

func GetBookings(c *gin.Context) {
	type BookingWithVendor struct {
		models.Booking
		Vendor models.VendorListing `json:"vendor"`
	}

	var bookings []BookingWithVendor

	query := `
		SELECT 
			b.id, b.user_id, b.vendor_id, b.event_date, b.status, b.notes, b.created_at, b.updated_at,
			v.id, v.vendor_id, v.title, v.description, v.category, v.price_range, v.location, v.photos, v.created_at, v.updated_at
		FROM bookings b
		JOIN vendors v ON b.vendor_id = v.id
		ORDER BY b.created_at DESC
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		log.Println("Query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b BookingWithVendor

		err := rows.Scan(
			// booking fields
			&b.ID,
			&b.UserID,
			&b.VendorID,
			&b.EventDate,
			&b.Status,
			&b.Notes,
			&b.CreatedAt,
			&b.UpdatedAt,
			// vendor fields
			&b.Vendor.ID,
			&b.Vendor.VendorID,
			&b.Vendor.Title,
			&b.Vendor.Description,
			&b.Vendor.Category,
			&b.Vendor.PriceRange,
			&b.Vendor.Location,
			&b.Vendor.Photos,
			&b.Vendor.CreatedAt,
			&b.Vendor.UpdatedAt,
		)

		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		bookings = append(bookings, b)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"bookings": bookings,
	})

}
