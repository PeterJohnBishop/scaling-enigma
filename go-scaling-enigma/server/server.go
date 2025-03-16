package server

import (
	"fmt"
	"log"
	"scaling-enigma/go-scaling-enigma/main.go/websocket"

	"github.com/gin-gonic/gin"
)

func ServeGin() {
	r := gin.Default()
	r.GET("/ws", websocket.ServeWebsocket)
	fmt.Println("WebSocket server running on ws://localhost:8080/ws")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Gin Server failed during startup: " + err.Error())
	}
}
