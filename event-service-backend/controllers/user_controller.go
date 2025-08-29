package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/dharmaseervi/event-service-backend/models"
	"github.com/dharmaseervi/event-service-backend/utils"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User

	log.Println("Creating user...")

	// Bind JSON payload to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("User data: %+v", user)

	// Insert user (no password now)
	query := `
		INSERT INTO users 
			(clerk_id, full_name, email, role, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at
	`

	err := config.DB.QueryRow(
		query,
		user.ClerkID,
		user.FullName,
		user.Email,
		user.Role,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not create user")
		return
	}

	log.Printf("User created with ID: %d", user.ID)
	utils.RespondWithJSON(c, http.StatusCreated, user)
}

func GetAllUsers(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, clerk_id, full_name, email, role, created_at, updated_at 
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.ClerkID,
			&user.FullName,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	utils.RespondWithJSON(c, http.StatusOK, users)
}
