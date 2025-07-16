package routes

import (
	"github.com/dharmaseervi/event-service-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.GET("/", controllers.GetAllUsers)
	}
}

func SetupVendorRoutes(router *gin.Engine) {
	vendorRoutes := router.Group("/vendors")
	{
		vendorRoutes.POST("/", controllers.CreateVendor)
		vendorRoutes.GET("/", controllers.GetAllVendors)
		vendorRoutes.GET("/:id", controllers.GetVendorByID)
		vendorRoutes.GET("/featured", controllers.GetFeaturedVendors)
		vendorRoutes.GET("/recommended", controllers.GetRecommendedVendors)
	}
}

func SetupSearchRoutes(router *gin.Engine) {
	searchRoutes := router.Group("/search")
	{
		searchRoutes.GET("/vendors", controllers.SearchVendors)
	}
}

func SetupSavedVendorRoutes(router *gin.Engine) {
	savedVendorRoutes := router.Group("/saved-vendors")
	{
		savedVendorRoutes.POST("/", controllers.SaveVendor)
		savedVendorRoutes.DELETE("/", controllers.UnsaveVendor)
		savedVendorRoutes.GET("/", controllers.GetSavedVendors)
	}
}

func SetupBookingRoutes(r *gin.Engine) {
	bookings := r.Group("/bookings")
	{
		bookings.POST("/", controllers.CreateBooking)
		bookings.GET("/", controllers.GetBookings)
	}
}

func SeTupSendNotification(r *gin.Engine) {
	notification := r.Group("/notif")
	{
		notification.POST("/", sendPushNotificationHandler)
		notification.POST("/push-token", sendPushNotificationHandler)
	}
}

// sendPushNotificationHandler wraps controllers.SendPushNotification to match gin.HandlerFunc
func sendPushNotificationHandler(c *gin.Context) {
	var req struct {
		PushToken string                 `json:"pushToken"`
		Title     string                 `json:"title"`
		Body      string                 `json:"body"`
		Data      map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	err := controllers.SendPushNotification(req.PushToken, req.Title, req.Body, req.Data)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Notification sent successfully"})
}
