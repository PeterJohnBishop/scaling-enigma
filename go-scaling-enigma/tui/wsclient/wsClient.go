package wsclient

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/gorilla/websocket"
)

const serverURL = "ws://localhost:8080/ws"

var (
	conn  *websocket.Conn
	mutex sync.Mutex
)

func ServeWebsocketClient() {
	var err error

	conn, _, err = websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("\rError connecting to WebSocket server:", err)
	}
	fmt.Println("\rTUI Connected to WebSocket server")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go listenForMessages()

	<-interrupt
	log.Println("Interrupt received, closing connection")

	if conn != nil {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Client is disconnecting"))
		if err != nil {
			log.Println("Close message error:", err)
		}
	}

	CloseConnection()
}

func listenForMessages() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("\rRead error:", err)
			return
		}
		fmt.Println("\rTUI Client Received:", string(message))
	}
}

func SendMessage(message string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if conn == nil {
		return fmt.Errorf("\rWebSocket connection is not established")
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("\rWrite error:", err)
		return err
	}
	return nil
}

func CloseConnection() {
	mutex.Lock()
	defer mutex.Unlock()

	if conn != nil {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("\rClose error:", err)
		}
		conn.Close()
		conn = nil
	}
}
