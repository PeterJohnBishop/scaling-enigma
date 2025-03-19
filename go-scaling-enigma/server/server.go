package server

import (
	"database/sql"
	"fmt"
	"os"

	"scaling-enigma/go-scaling-enigma/main.go/server/auth"
	"scaling-enigma/go-scaling-enigma/main.go/server/routes"
	"scaling-enigma/go-scaling-enigma/main.go/server/websocket"

	"github.com/gin-gonic/gin"
)

func ServeGin(db *sql.DB) {

	auth.Init()

	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = f

	r := gin.Default()
	addDefaultRoutes(r)
	addUserRoutes(db, r)
	go websocket.HandleBroadcast()
	fmt.Println("Server listening on [ Gin http://localhost:8080 ] : [ WebSocket ws://localhost:8080/ws ]")
	r.Run(":8080")
}

func addDefaultRoutes(r *gin.Engine) {
	r.GET("/", routes.HelloHandler)
	r.GET("/ws", websocket.ServeWebsocket)
}

func addUserRoutes(db *sql.DB, r *gin.Engine) {
	r.POST("/users/new", func(c *gin.Context) {
		routes.CreateUserHandler(db, c)
	})
	r.POST("/users/login", func(c *gin.Context) {
		routes.Login(db, c)
	})
	r.GET("/users/all", func(c *gin.Context) {
		routes.GetUsersHandler(db, c)
	})
	r.GET("/users/user/:id", func(c *gin.Context) {
		routes.GetUserByIDHandler(db, c)
	})
	r.PUT("/users/user/update", func(c *gin.Context) {
		routes.UpdateUserHandler(db, c)
	})
	r.DELETE("/users/user/delete/:id", func(c *gin.Context) {
		routes.DeleteUserHandler(db, c)
	})
}
