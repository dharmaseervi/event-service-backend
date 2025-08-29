package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/gin-gonic/gin"
)

// POST /bookings
func CreateBooking(c *gin.Context) {
	var booking models.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		log.Printf("BindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Clerk middleware verified the JWT and put the session claims in the request context
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	clerkID := claims.Subject // <- This is the Clerk user ID (e.g. "user_abc123")

	fmt.Println("Clerk ID:", clerkID)

	// Map Clerk ID -> local DB user id
	var localUserID int
	if err := config.DB.QueryRow(`SELECT id FROM users WHERE clerk_id=$1`, clerkID).Scan(&localUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	booking.UserID = localUserID

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

	log.Printf("❌ Insert failed: %v", err)

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

func GetMyBookings(c *gin.Context) {
	type BookingWithVendor struct {
		models.Booking
		Vendor models.VendorListing `json:"vendor"`
	}

	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	clerkID := claims.Subject

	var localUserID int
	if err := config.DB.QueryRow(`SELECT id FROM users WHERE clerk_id=$1`, clerkID).Scan(&localUserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	const q = `
		SELECT
		  b.id, b.user_id, b.vendor_id, b.event_date, b.status, b.notes, b.created_at, b.updated_at,
		  v.id, v.vendor_id, v.title, v.description, v.category, v.price_range, v.location,
		  COALESCE(v.photos, ARRAY[]::text[]) AS photos,
		  v.created_at, v.updated_at
		FROM bookings b
		LEFT JOIN vendors v ON v.id = b.vendor_id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`
	rows, err := config.DB.Query(q, localUserID)
	if err != nil {
		log.Printf("GetMyBookings query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}
	defer rows.Close()

	var out []BookingWithVendor
	count := 0
	for rows.Next() {
		var b BookingWithVendor
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.VendorID, &b.EventDate, &b.Status, &b.Notes, &b.CreatedAt, &b.UpdatedAt,
			&b.Vendor.ID, &b.Vendor.VendorID, &b.Vendor.Title, &b.Vendor.Description, &b.Vendor.Category,
			&b.Vendor.PriceRange, &b.Vendor.Location, &b.Vendor.Photos,
			&b.Vendor.CreatedAt, &b.Vendor.UpdatedAt,
		); err != nil {
			log.Printf("GetMyBookings scan error: %v", err) // ← SEE THIS IN LOGS
			continue
		}
		out = append(out, b)
		count++
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetMyBookings rows iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "iteration failed"})
		return
	}

	log.Printf("GetMyBookings: user %d → %d rows", localUserID, count)
	c.JSON(http.StatusOK, gin.H{"success": true, "bookings": out})
}
