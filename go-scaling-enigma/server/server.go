package server

import (
	"fmt"
	"log"
	"net/http"
	"scaling-enigma/go-scaling-enigma/main.go/websocket"

	"github.com/gin-gonic/gin"
)

func ServeGin() {
	r := gin.Default()
	r.GET("/", helloHandler)
	r.GET("/ws", websocket.ServeWebsocket)
	fmt.Println("WebSocket server running on ws://localhost:8080/ws")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Gin Server failed during startup: " + err.Error())
	}
}

func helloHandler(c *gin.Context) {
	message := map[string]string{"message": "Hello"}
	c.IndentedJSON(http.StatusOK, message)
}
