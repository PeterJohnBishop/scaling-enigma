package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Failed to upgrade server to Websocket Protocol: " + err.Error())
	}
	defer conn.Close()

	var wg sync.WaitGroup

	// send event
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()

	// listen for events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				break
			}
			fmt.Println("Received:", string(msg))
		}
	}()

	select {}
}
