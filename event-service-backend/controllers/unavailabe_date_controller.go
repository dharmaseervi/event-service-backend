package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/dharmaseervi/event-service-backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateUnavailableDate(c *gin.Context) {
	var date models.UnavailableDate

	log.Println("📥 Incoming request body for unavailable date...")

	if err := c.ShouldBindJSON(&date); err != nil {
		log.Println("❌ Error binding JSON:", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}

	vendorIDStr := c.Param("vendor_id")
	if vendorIDStr == "" {
		log.Println("❌ Missing vendor_id in URL")
		utils.RespondWithError(c, http.StatusBadRequest, "Vendor ID is required")
		return
	}

	var vendorID int
	if _, err := fmt.Sscanf(vendorIDStr, "%d", &vendorID); err != nil {
		log.Println("❌ Invalid vendor_id format:", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid vendor ID format")
		return
	}

	log.Printf("📌 Parsed vendor_id: %d | Date: %+v\n", vendorID, date)

	query := `
		INSERT INTO vendor_bookings (vendor_id, booked_from, booked_to)
		VALUES ($1, $2, $3)
		RETURNING id, booked_from, booked_to
	`

	err := config.DB.QueryRow(
		query,
		vendorID,
		date.BookedFrom,
		date.BookedTo,
	).Scan(&date.ID, &date.BookedFrom, &date.BookedTo)

	if err != nil {
		log.Printf("❌ DB error inserting unavailable date: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not create unavailable date")
		return
	}

	log.Println("✅ Unavailable date inserted successfully")
	c.JSON(http.StatusCreated, gin.H{
		"message":        "Unavailable date created successfully",
		"unavailable_id": date.ID,
		"booked_from":    date.BookedFrom,
		"booked_to":      date.BookedTo,
	})
}

func GetVendorUnavailability(c *gin.Context) {
	vendorID := c.Param("vendor_id")
	if vendorID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Vendor ID is required")
		return
	}

	query := `
		SELECT booked_from, booked_to 
		FROM vendor_bookings 
		WHERE vendor_id = $1
	`
	rows, err := config.DB.Query(query, vendorID)
	if err != nil {
		log.Printf("Error fetching unavailability: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not fetch unavailability")
		return
	}
	defer rows.Close()

	var unavailability []models.UnavailableDate
	for rows.Next() {
		var date models.UnavailableDate
		if err := rows.Scan(&date.BookedFrom, &date.BookedTo); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		unavailability = append(unavailability, date)
	}

	c.JSON(http.StatusOK, unavailability)
}
