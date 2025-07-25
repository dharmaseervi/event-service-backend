package controllers

import (
	"fmt"
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
	searchQuery := strings.TrimSpace(c.Query("q"))
	location := strings.TrimSpace(c.Query("location"))
	category := strings.TrimSpace(c.Query("category"))
	fromDate := strings.TrimSpace(c.Query("from_date"))
	toDate := strings.TrimSpace(c.Query("to_date"))

	if searchQuery == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Search query is required")
		return
	}

	vendors := []models.VendorListing{}

	query := `
	SELECT 
		id, vendor_id, title, description, category, 
		price_range, location, photos, created_at, updated_at,
		ts_rank(search_vector, websearch_to_tsquery('english', $1)) as rank
	FROM vendors
	WHERE search_vector @@ websearch_to_tsquery('english', $1)
	`
	args := []interface{ any }{searchQuery}
	argPos := 2

	if location != "" {
		query += fmt.Sprintf(" AND location ILIKE $%d", argPos)
		args = append(args, "%"+location+"%")
		argPos++
	}

	if category != "" {
		query += fmt.Sprintf(" AND category ILIKE $%d", argPos)
		args = append(args, "%"+category+"%")
		argPos++
	}

	if fromDate != "" && toDate != "" {
		query += fmt.Sprintf(`
			AND NOT EXISTS (
				SELECT 1 FROM vendor_bookings vb
				WHERE vb.vendor_id = vendors.id
				AND vb.booked_from <= $%d
				AND vb.booked_to >= $%d
			)
		`, argPos, argPos+1)
		args = append(args, toDate, fromDate)
		argPos += 2
	}

	query += " ORDER BY rank DESC, created_at DESC LIMIT 50"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error searching vendors: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not perform search")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var vendor models.VendorListing
		var photos pq.StringArray
		var loc string
		var rank float64

		err := rows.Scan(
			&vendor.ID,
			&vendor.VendorID,
			&vendor.Title,
			&vendor.Description,
			&vendor.Category,
			&vendor.PriceRange,
			&loc,
			&photos,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
			&rank,
		)
		if err != nil {
			log.Printf("Error scanning result: %v", err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to read results")
			return
		}
		vendor.Location = loc
		vendor.Photos = []string(photos)
		vendors = append(vendors, vendor)
	}

	utils.RespondWithJSON(c, http.StatusOK, vendors)
}
