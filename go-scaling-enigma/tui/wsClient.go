package tui

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func ServeWebsocketClient() {

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Websocket Client failed to connect: " + err.Error())
	}
	defer conn.Close()

	var wg sync.WaitGroup

	// listen for events
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("\rRead error:", err)
				break
			}
			fmt.Println("\rReceived:", string(message))
		}
	}()

	// send event
	wg.Add(1)
	go func() {
		defer wg.Done()
		go func() {
			for {
				time.Sleep(3 * time.Second)
				msg := "\rcli_event: Hello from the CLI!"
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					fmt.Println("\rWrite error:", err)
					break
				}
			}
		}()
	}()

	select {} // keeps connection running

}
