package controllers

import (
	"net/http"
	"strconv"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateVendorDeals(c *gin.Context) {
	var deal models.VendorDeal
	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vendorID, err := strconv.Atoi(c.Param("vendor_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vendor_id"})
		return
	}
	deal.VendorID = vendorID

	query := `
			INSERT INTO vendor_deals (vendor_id, title, description, discount_percent, original_price, deal_price, start_date, end_date, photos)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, created_at
		`
	err = config.DB.QueryRow(
		query,
		deal.VendorID, deal.Title, deal.Description, deal.DiscountPercent,
		deal.OriginalPrice, deal.DealPrice, deal.StartDate, deal.EndDate, pq.Array(deal.Photos),
	).Scan(&deal.ID, &deal.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, deal)
}

func GetAllVendorsDeals(c *gin.Context) {

	rows, err := config.DB.Query(`
			SELECT id, vendor_id, title, description, discount_percent, original_price, deal_price, start_date, end_date, photos, created_at
			FROM vendor_deals 
		`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var deals []models.VendorDeal
	for rows.Next() {
		var deal models.VendorDeal
		var photos []string
		err := rows.Scan(&deal.ID, &deal.VendorID, &deal.Title, &deal.Description, &deal.DiscountPercent,
			&deal.OriginalPrice, &deal.DealPrice, &deal.StartDate, &deal.EndDate, pq.Array(&photos), &deal.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		deal.Photos = photos
		deals = append(deals, deal)
	}
	c.JSON(http.StatusOK, deals)
}
