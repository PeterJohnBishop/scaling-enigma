package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ServeGin() {

	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = f

	r := gin.Default()
	r.GET("/", helloHandler)
	r.GET("/ws", ServeWebsocket)

	log.Println("Server: [ Gin http://localhost:8080 ] : [ WebSocket ws://localhost:8080/ws ]")
	r.Run(":8080")
}

func helloHandler(c *gin.Context) {
	message := map[string]string{"message": "Hello"}
	c.IndentedJSON(http.StatusOK, message)
}
