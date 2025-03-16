package websocket

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func WebsocketClient() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	defer conn.Close()

	// listen for events
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Received:", string(message))
	}

	// send event
	go func() {
		for {
			time.Sleep(3 * time.Second)
			msg := "cli_event: Hello from the CLI!"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println("Write error:", err)
				break
			}
		}
	}()

	select {} // keeps connection running

}
