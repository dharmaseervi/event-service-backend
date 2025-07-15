package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/migrations"
	"github.com/dharmaseervi/event-service-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()
	defer config.CloseDB()
	// Run migrations
	if err := migrations.RunMigrations(config.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupUserRoutes(router)
	routes.SetupVendorRoutes(router)
	routes.SetupSearchRoutes(router)
	routes.SetupSavedVendorRoutes(router)
	routes.SetupBookingRoutes(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server is running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Block until we receive a signal
	<-quit
	log.Println("Shutting down server...")

	// Close database connection
	config.CloseDB()
	log.Println("Server gracefully stopped")
}
