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

func CreateChatHandler(db *sql.DB, c *gin.Context) {

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

	var newChat postgres.Chat
	if err := c.ShouldBindJSON(&newChat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	chat, err := postgres.CreateChat(db, newChat)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	event := websocket.Event{
		Type:    "UserCreated",
		Message: fmt.Sprintf("Chat created successfully: %s", chat),
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusCreated, chat)
}

func GetChatsHandler(db *sql.DB, c *gin.Context) {
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
	var chats []postgres.Chat
	allChats, err := postgres.GetChats(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all chats"})
		return
	}
	chats = allChats

	event := websocket.Event{
		Type:    "ChatsFound",
		Message: fmt.Sprintf("%d chats found", len(chats)),
	}
	websocket.BroadcastMessage(event)
	c.JSON(http.StatusOK, chats)
}

func UpdateChatHandler(db *sql.DB, c *gin.Context) {
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
	var chat postgres.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedChat, err := postgres.UpdateChatByID(db, chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	event := websocket.Event{
		Type:    "ChatUpdated",
		Message: fmt.Sprintf("User updated: %s", updatedChat),
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusOK, gin.H{
		"message": "Chat updated successfully",
		"user":    updatedChat,
	})
}

func DeleteChatHandler(db *sql.DB, c *gin.Context) {
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
	err := postgres.DeleteChatById(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	event := websocket.Event{
		Type:    "ChatDeletd",
		Message: "Chat Deleted",
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
