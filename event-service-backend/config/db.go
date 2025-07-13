package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Build connection string from environment variables
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_SSLMODE"),
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	log.Println("Successfully connected to database!")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
