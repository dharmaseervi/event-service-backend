package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dharmaseervi/event-service-backend/config"
	"github.com/gin-gonic/gin"
)

type ExpoMessage struct {
	To    string                 `json:"to"`
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Sound string                 `json:"sound,omitempty"`
	Data  map[string]interface{} `json:"data,omitempty"`
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
		c.JSON(500, gin.H{"error": "Failed to store token"})
		return
	}

	c.JSON(200, gin.H{"message": "Token stored"})
}

func SendPushNotification(pushToken, title, body string, data map[string]interface{}) error {
	message := ExpoMessage{
		To:    pushToken,
		Title: title,
		Body:  body,
		Sound: "default",
		Data:  data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://exp.host/--/api/v2/push/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Expo push response status: %s", resp.Status)
	return nil
}
