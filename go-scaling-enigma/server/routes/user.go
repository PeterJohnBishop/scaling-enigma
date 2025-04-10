package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"scaling-enigma/go-scaling-enigma/main.go/server/auth"
	"scaling-enigma/go-scaling-enigma/main.go/server/postgres"
	"scaling-enigma/go-scaling-enigma/main.go/server/websocket"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func HelloHandler(c *gin.Context) {
	message := map[string]string{"message": "Hello"}
	c.IndentedJSON(http.StatusOK, message)
}

func CreateUserHandler(db *sql.DB, c *gin.Context) {
	var newUser postgres.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := postgres.CreateUser(db, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	event := websocket.Event{
		Type:    "UserCreated",
		Message: fmt.Sprintf("User created successfully: %s", user),
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *sql.DB, c *gin.Context) {

	var req LoginRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode request body"})
		return
	}

	user, err := postgres.GetUserByEmail(db, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by that email"})
		return
	}
	pass := auth.CheckPasswordHash(req.Password, user.Password)
	if !pass {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password Verfication Failed"})
		return
	}

	userClaims := auth.UserClaims{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}

	token, err := auth.NewAccessToken(userClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	refreshToken, err := auth.NewRefreshToken(userClaims.StandardClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	event := websocket.Event{
		Type:    "UserLogin",
		Message: fmt.Sprintf("User logged in successfully: %s", user),
	}
	websocket.BroadcastMessage(event)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login Success",
		"token":         token,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func GetUserByEmailHandler(db *sql.DB, email string, c *gin.Context) {
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
	parseErr := auth.ParseAccessToken(token)
	if parseErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token!"})
		return
	}
	var user postgres.User
	foundUser, err := postgres.GetUserByEmail(db, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by that email"})
		return
	}

	user = foundUser

	event := websocket.Event{
		Type:    "UserFoundByEmail",
		Message: fmt.Sprintf("User found by email: %s", user),
	}
	websocket.BroadcastMessage(event)
	c.JSON(http.StatusOK, user)
}

func GetUserByIDHandler(db *sql.DB, c *gin.Context) {
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
	var user postgres.User
	id := c.Param("id")
	foundUser, err := postgres.GetUserByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by that id"})
		return
	}

	user = foundUser

	event := websocket.Event{
		Type:    "UserFoundById",
		Message: fmt.Sprintf("User found by id: %s", user),
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusOK, user)
}

func GetUsersHandler(db *sql.DB, c *gin.Context) {
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
	var users []postgres.User
	allUsers, err := postgres.GetUsers(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all users"})
		return
	}
	users = allUsers

	event := websocket.Event{
		Type:    "UsersFound",
		Message: fmt.Sprintf("%d users found", len(users)),
	}
	websocket.BroadcastMessage(event)
	c.JSON(http.StatusOK, users)
}

func UpdateUserHandler(db *sql.DB, c *gin.Context) {
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
	var user postgres.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedUser, err := postgres.UpdateUserByID(db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	event := websocket.Event{
		Type:    "UserUpdated",
		Message: fmt.Sprintf("User updated: %s", updatedUser),
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}

func DeleteUserHandler(db *sql.DB, c *gin.Context) {
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
	err := postgres.DeleteUserByID(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	event := websocket.Event{
		Type:    "UserDeleted",
		Message: "User Deleted",
	}
	websocket.BroadcastMessage(event)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
