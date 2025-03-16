package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close()

	// send event
	go func() {
		for {
			time.Sleep(5 * time.Second)
			msg := "server_event: Hello from the server!"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println("Write error:", err)
				break
			}
		}
	}()

	// listen for events
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Received:", string(msg))
	}
}

func OpenWebsocket() {
	r := gin.Default()
	r.GET("/ws", serveWs)
	fmt.Println("WebSocket server running on ws://localhost:8080/ws")
	r.Run(":8080")
}
