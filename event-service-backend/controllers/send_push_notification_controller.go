package controllers

import (
	"log"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/gin-gonic/gin"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

type ExpoMessage struct {
	To    string            `json:"to"`
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Sound string            `json:"sound,omitempty"`
	Data  map[string]string `json:"data"`
}

func SavePushToken(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}

	// Bind only the token from client
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 🔑 Clerk ID from middleware
	clerkID := c.GetString("clerk_id")
	if clerkID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	// 🔍 Map Clerk ID → local user id
	var localUserID int
	if err := config.DB.QueryRow(`SELECT id FROM users WHERE clerk_id=$1`, clerkID).Scan(&localUserID); err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	// 💾 Store push token tied to local user id
	_, err := config.DB.Exec(`
		INSERT INTO push_tokens (user_id, token) 
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET token=$2, updated_at=CURRENT_TIMESTAMP
	`, localUserID, req.Token)

	if err != nil {
		log.Printf("DB error saving push token: %v", err)
		c.JSON(500, gin.H{"error": "Failed to store token"})
		return
	}

	c.JSON(200, gin.H{"message": "Token stored"})
}

func SendPushNotification(c *gin.Context) {
	var req struct {
		UserID int               `json:"user_id"`
		Title  string            `json:"title"`
		Body   string            `json:"body"`
		Sound  string            `json:"sound,omitempty"`
		Data   map[string]string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 🔍 Get push token from DB
	var token string
	err := config.DB.QueryRow(`SELECT token FROM push_tokens WHERE user_id=$1`, req.UserID).Scan(&token)
	if err != nil {
		c.JSON(404, gin.H{"error": "No push token for this user"})
		return
	}

	pushToken, err := expo.NewExponentPushToken(token)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Expo push token"})
		return
	}

	client := expo.NewPushClient(nil)

	msg := &expo.PushMessage{
		To:       []expo.ExponentPushToken{pushToken},
		Title:    req.Title,
		Body:     req.Body,
		Sound:    req.Sound,
		Data:     req.Data,
		Priority: expo.DefaultPriority,
	}

	resp, err := client.Publish(msg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to publish notification"})
		return
	}

	if err := resp.ValidateResponse(); err != nil {
		c.JSON(500, gin.H{"error": "Expo rejected message"})
		return
	}

	c.JSON(200, gin.H{"success": true, "response": resp})
}
