package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/dharmaseervi/event-service-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateVendor(c *gin.Context) {
	var vendor models.VendorListing
	if err := c.ShouldBindJSON(&vendor); err != nil {
		log.Printf("Error binding JSON: %v", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input")
		return
	}
	log.Printf("Vendor data: %+v", vendor.VendorID)

	query := `
		INSERT INTO vendors 
			(vendor_id, title, description, category, price_range, location, photos, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := config.DB.QueryRow(
		query,
		vendor.VendorID,
		vendor.Title,
		vendor.Description,
		vendor.Category,
		vendor.PriceRange,
		vendor.Location,
		pq.StringArray(vendor.Photos), // ✅ convert []string properly
		time.Now(),
		time.Now(),
	).Scan(&vendor.ID, &vendor.CreatedAt, &vendor.UpdatedAt)

	if err != nil {
		log.Printf("Error creating vendor: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not create vendor")
		return
	}

	c.JSON(http.StatusCreated, vendor)
}

func GetAllVendors(c *gin.Context) {
	// Get category filter from query params
	category := c.Query("category")

	var query string
	var rows *sql.Rows
	var err error

	if category != "" {
		query = `
			SELECT id, vendor_id, title, description, category, price_range, location, photos, created_at, updated_at
			FROM vendors 
			WHERE category = $1
			ORDER BY created_at DESC
		`
		rows, err = config.DB.Query(query, category)
	} else {
		query = `
			SELECT id, vendor_id, title, description, category, price_range, location, photos, created_at, updated_at
			FROM vendors
			ORDER BY created_at DESC
		`
		rows, err = config.DB.Query(query)
	}

	if err != nil {
		log.Printf("Error querying vendors: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not retrieve vendors")
		return
	}
	defer rows.Close()

	vendors := []models.VendorListing{}
	for rows.Next() {
		var vendor models.VendorListing
		if err := rows.Scan(
			&vendor.ID,
			&vendor.VendorID,
			&vendor.Title,
			&vendor.Description,
			&vendor.Category,
			&vendor.PriceRange,
			&vendor.Location,
			(*pq.StringArray)(&vendor.Photos),
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning vendor: %v", err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Could not retrieve vendors")
			return
		}
		vendors = append(vendors, vendor)
	}

	utils.RespondWithJSON(c, http.StatusOK, vendors)
}

func GetVendorByID(c *gin.Context) {
	id := c.Param("id") // Get ID from URL
	log.Printf("Fetching vendor with ID: %s", id)
	var vendor models.VendorListing

	query := `
		SELECT id, vendor_id, title, description, category, price_range, location, photos, created_at, updated_at
		FROM vendors WHERE id = $1
	`
	err := config.DB.QueryRow(query, id).Scan(
		&vendor.ID,
		&vendor.VendorID,
		&vendor.Title,
		&vendor.Description,
		&vendor.Category,
		&vendor.PriceRange,
		&vendor.Location,
		(*pq.StringArray)(&vendor.Photos),
		&vendor.CreatedAt,
		&vendor.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error fetching vendor by ID: %v", err)
		utils.RespondWithError(c, http.StatusNotFound, "Vendor not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, vendor)
}

func DeleteVendor(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Deleting vendor with ID: %s", id)

	query := `DELETE FROM vendors WHERE id = $1`

	result, err := config.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting vendor: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not delete vendor")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.RespondWithError(c, http.StatusNotFound, "Vendor not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor deleted successfully"})
}
