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
