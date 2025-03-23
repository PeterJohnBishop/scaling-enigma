package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"scaling-enigma/go-scaling-enigma/main.go/server/auth"
	"scaling-enigma/go-scaling-enigma/main.go/server/postgres"
	"scaling-enigma/go-scaling-enigma/main.go/server/websocket"

	"github.com/gin-gonic/gin"
)

func CreateMessageHandler(db *sql.DB, c *gin.Context) {

	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token missing!"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format!"})
		return
	}
	claims := auth.ParseAccessToken(token)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token!"})
		return
	}

	var newMessage postgres.Message
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	message, err := postgres.CreateMessage(db, newMessage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	event := websocket.Event{
		Type:    "MessageCreated",
		Message: fmt.Sprintf("Message created successfully: %s", message),
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusCreated, message)
}

func GetMessagesHandler(db *sql.DB, c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token missing!"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format!"})
		return
	}
	claims := auth.ParseAccessToken(token)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token!"})
		return
	}
	var messages []postgres.Message
	allMessages, err := postgres.GetMessages(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all chats"})
		return
	}
	messages = allMessages

	event := websocket.Event{
		Type:    "MessagesFound",
		Message: fmt.Sprintf("%d messages found", len(messages)),
	}
	websocket.BroadcastMessage(event)
	c.JSON(http.StatusOK, messages)
}

func DeleteMessageHandler(db *sql.DB, c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token missing!"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format!"})
		return
	}
	claims := auth.ParseAccessToken(token)
	if claims == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token!"})
		return
	}
	id := c.Param("id")
	err := postgres.DeleteMessageById(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	event := websocket.Event{
		Type:    "MessageDeleted",
		Message: "Message Deleted",
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
