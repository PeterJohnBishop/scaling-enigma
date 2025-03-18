package server

import (
	"database/sql"
	"log"
	"os"

	"scaling-enigma/go-scaling-enigma/main.go/server/auth"
	"scaling-enigma/go-scaling-enigma/main.go/server/routes"

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

	log.Println("Server: [ Gin http://localhost:8080 ] : [ WebSocket ws://localhost:8080/ws ]")
	r.Run(":8080")
}

func addDefaultRoutes(r *gin.Engine) {
	r.GET("/", routes.HelloHandler)
	r.GET("/ws", ServeWebsocket)
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
	r.DELETE("/users/user/:id", func(c *gin.Context) {
		routes.DeleteUserHandler(db, c)
	})
}
