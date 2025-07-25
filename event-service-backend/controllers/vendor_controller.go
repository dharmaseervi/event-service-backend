package controllers

import (
	"database/sql"
	"fmt"
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
	log.Printf("Vendor data: %+v", vendor)

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
	category := c.Param("category")

	log.Printf("Fetching vendor with ID: %s", id)
	var vendor models.VendorListing

	query := `
		SELECT id, vendor_id, title, description, category, price_range, location, photos, created_at, updated_at
		FROM vendors WHERE id = $1
	`
	if category != "" {
		query += ` AND category = $2`
	}

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

// GET /vendors/featured
func GetFeaturedVendors(c *gin.Context) {
	var featured []models.VendorListing

	query := `SELECT id, vendor_id, title, description, category, price_range, location, photos, rating, featured, created_at, updated_at 
FROM vendors WHERE featured = true ORDER BY updated_at DESC`

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch featured vendors"})
		return
	}
	defer rows.Close()

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
			&vendor.Photos,
			&vendor.Rating,
			&vendor.Featured,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
		); err != nil {
			log.Println("❌ Failed to scan vendor:", err)
			continue
		}
		featured = append(featured, vendor)
	}

	c.JSON(http.StatusOK, gin.H{"vendors": featured})
}

func GetRecommendedVendors(c *gin.Context) {
	location := c.Query("location") // Optional
	category := c.Query("category") // Optional

	var recommendations []models.VendorListing

	query := `
		SELECT id, title, category, description, location, price_range, photos, rating, featured
		FROM vendors
		WHERE featured = true
	`

	var args []any
	argIndex := 1

	if category != "" {
		query += ` AND category = $` + fmt.Sprint(argIndex)
		args = append(args, category)
		argIndex++
	}

	if location != "" {
		query += ` AND location ILIKE '%' || $` + fmt.Sprint(argIndex) + ` || '%'`
		args = append(args, location)
		argIndex++
	}

	query += ` ORDER BY rating DESC, updated_at DESC LIMIT 10`

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Printf("Failed to fetch recommendations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var vendor models.VendorListing
		err := rows.Scan(&vendor.ID, &vendor.Title, &vendor.Category, &vendor.Description, &vendor.Location, &vendor.PriceRange, &vendor.Photos, &vendor.Rating, &vendor.Featured)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		recommendations = append(recommendations, vendor)
	}

	c.JSON(http.StatusOK, gin.H{"vendors": recommendations})
}
