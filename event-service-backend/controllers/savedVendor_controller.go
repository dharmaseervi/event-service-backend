package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/gin-gonic/gin"
)

// POST /saved-vendors
// body: { "vendor_id": number }
func SaveVendor(c *gin.Context) {
	// 1) Auth: clerk -> local user id
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	clerkID := claims.Subject

	print("Clerk ID: ", clerkID, "\n")
	var localUserID int64
	if err := config.DB.QueryRow(`SELECT id FROM users WHERE clerk_id=$1`, clerkID).Scan(&localUserID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	print("local ID: ", localUserID, "\n")
	// 2) Parse payload (only vendor_id)
	var input struct {
		VendorID int64 `json:"vendor_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.VendorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vendor_id is required"})
		return
	}

	fmt.Println("VendorID:", input.VendorID)
	// (Optional) validate vendor exists
	var exists bool
	if err := config.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM vendors WHERE id=$1)`, input.VendorID).Scan(&exists); err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vendor_id"})
		return
	}

	// 3) Insert (idempotent thanks to UNIQUE(user_id, vendor_id))
	var id int64
	var createdAt time.Time

	err := config.DB.QueryRow(`
		INSERT INTO saved_items (user_id, vendor_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, vendor_id) DO UPDATE SET created_at = EXCLUDED.created_at
		RETURNING id, created_at
	`, localUserID, input.VendorID, time.Now()).Scan(&id, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save vendor"})
		fmt.Println("Error inserting saved vendor:", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"user_id":    localUserID,
		"vendor_id":  input.VendorID,
		"created_at": createdAt,
	})
}

// DELETE /saved-vendors/:vendor_id
func UnsaveVendor(c *gin.Context) {
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	clerkID := claims.Subject

	var localUserID int64
	if err := config.DB.QueryRow(`SELECT id FROM users WHERE clerk_id=$1`, clerkID).Scan(&localUserID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	vendorIDStr := c.Param("vendor_id")
	var vendorID int64
	if _, err := fmt.Sscanf(vendorIDStr, "%d", &vendorID); err != nil || vendorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vendor_id"})
		return
	}

	if _, err := config.DB.Exec(`DELETE FROM saved_items WHERE user_id=$1 AND vendor_id=$2`, localUserID, vendorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not unsave vendor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vendor unsaved"})
}

// GET /saved-vendors/me
func GetMySavedVendors(c *gin.Context) {
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	clerkID := claims.Subject

	var localUserID int64
	if err := config.DB.QueryRow(`SELECT id FROM users WHERE clerk_id=$1`, clerkID).Scan(&localUserID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	rows, err := config.DB.Query(`
		SELECT v.id, v.vendor_id, v.title, v.description, v.category, v.price_range,
		       v.location, v.photos, v.rating, v.featured, v.created_at, v.updated_at
		FROM saved_items sv
		JOIN vendors v ON sv.vendor_id = v.id
		WHERE sv.user_id = $1
		ORDER BY sv.created_at DESC
	`, localUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var vendors []models.VendorListing
	for rows.Next() {
		var v models.VendorListing
		if err := rows.Scan(
			&v.ID, &v.VendorID, &v.Title, &v.Description, &v.Category, &v.PriceRange,
			&v.Location, &v.Photos, &v.Rating, &v.Featured, &v.CreatedAt, &v.UpdatedAt,
		); err == nil {
			vendors = append(vendors, v)
		}
	}

	c.JSON(http.StatusOK, gin.H{"vendors": vendors})
}
