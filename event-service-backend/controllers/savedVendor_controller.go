package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/dharmaseervi/event-service-backend/utils"

	"github.com/gin-gonic/gin"
)

// Save a vendor for the user
func SaveVendor(c *gin.Context) {

	var input models.SavedVendor

	log.Printf("Received request to save vendor: %v", c.Request.Body)
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	log.Printf("Saving vendor for user: %s vendor: %d", input.UserID, input.VendorID)

	if input.UserID == "" || input.VendorID == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "Missing user_id or vendor_id")
		return
	}

	query := `
		INSERT INTO saved_vendors (user_id, vendor_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := config.DB.QueryRow(query, input.UserID, input.VendorID, time.Now()).Scan(&input.ID, &input.CreatedAt)
	if err != nil {
		log.Printf("SQL insert error: %v", err) // 👈 log exact DB error
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not save vendor")
		return
	}

	c.JSON(http.StatusCreated, input)
}

// Remove a saved vendor
func UnsaveVendor(c *gin.Context) {
	userID := c.Query("user_id")
	vendorID := c.Query("vendor_id")

	if userID == "" || vendorID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Missing user_id or vendor_id")
		return
	}

	query := `DELETE FROM saved_vendors WHERE user_id = $1 AND vendor_id = $2`
	_, err := config.DB.Exec(query, userID, vendorID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not unsave vendor")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor unsaved"})
}

// Get all saved vendors for a user
func GetSavedVendors(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Missing user_id")
		return
	}

	query := `SELECT id, user_id, vendor_id, created_at FROM saved_vendors WHERE user_id = $1`
	rows, err := config.DB.Query(query, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch saved vendors")
		return
	}
	defer rows.Close()

	var saved []models.SavedVendor

	for rows.Next() {
		var s models.SavedVendor
		if err := rows.Scan(&s.ID, &s.UserID, &s.VendorID, &s.CreatedAt); err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Scan error")
			return
		}
		saved = append(saved, s)
	}

	log.Printf("Found %d saved vendors for user %s", len(saved), userID)
	c.JSON(http.StatusOK, saved)
}
