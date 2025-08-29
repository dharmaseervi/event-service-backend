package controllers

import (
	"fmt"
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
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	_, err := config.DB.Exec(`INSERT INTO push_tokens (user_id, token) VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET token = $2, updated_at = CURRENT_TIMESTAMP`, req.UserID, req.Token)

	if err != nil {
		log.Printf("hey: %v", err)
		c.JSON(500, gin.H{"error": "Failed to store token"})
		return
	}

	c.JSON(200, gin.H{"message": "Token stored"})
}

func SendPushNotification(c *gin.Context) {
	var req ExpoMessage

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	pushToken, err := expo.NewExponentPushToken("ExponentPushToken[81NJ4oBLW2SpwlmCXMFCby]")
	if err != nil {
		panic(err)
	}

	// Create a new Expo SDK client
	client := expo.NewPushClient(nil)

	// Publish message
	response, err := client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     req.Body,
			Data:     req.Data,
			Sound:    "default",
			Title:    req.Title,
			Priority: expo.DefaultPriority,
		},
	)

	// Check errors
	if err != nil {
		panic(err)
	}

	// Validate responses
	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
	}

	fmt.Println([]byte(pushToken))

}
