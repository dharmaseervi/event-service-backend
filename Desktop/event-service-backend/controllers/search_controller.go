package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/dharmaseervi/event-service-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func SearchVendors(c *gin.Context) {
	// Get and validate search query
	searchQuery := strings.TrimSpace(c.Query("q"))
	log.Printf("Search query: %s", searchQuery)
	if searchQuery == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Search query is required")
		return
	}

	vendors := []models.VendorListing{}

	// Full-text search query with ranking and sanitization
	query := `
	SELECT 
		id, vendor_id, title, description, category, 
		price_range, location, photos, created_at, updated_at,
		ts_rank(search_vector, websearch_to_tsquery('english', $1)) as rank
	FROM vendors
	WHERE search_vector @@ websearch_to_tsquery('english', $1)
	ORDER BY rank DESC, created_at DESC
	LIMIT 50
`

	rows, err := config.DB.Query(query, searchQuery)
	if err != nil {
		log.Printf("Error searching vendors: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not perform search")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var vendor models.VendorListing
		var photos pq.StringArray
		var location string
		var rank float64

		err := rows.Scan(
			&vendor.ID,
			&vendor.VendorID,
			&vendor.Title,
			&vendor.Description,
			&vendor.Category,
			&vendor.PriceRange,
			&location, // first location (string)
			&photos,   // then photos (pq.StringArray)
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
			&rank,
		)

		vendor.Location = location
		vendor.Photos = []string(photos)

		if err != nil {
			log.Printf("Error scanning search result: %v", err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Could not process results")
			return
		}
		vendors = append(vendors, vendor)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after scanning rows: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Error processing results")
		return
	}

	// Return results with search metadata
	utils.RespondWithJSON(c, http.StatusOK, vendors)
}
