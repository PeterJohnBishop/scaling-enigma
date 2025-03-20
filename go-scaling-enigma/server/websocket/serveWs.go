package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Event struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Event)
	mutex     sync.Mutex
)

func ServeWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("\rFailed to upgrade connection:", err)
		return
	}

	defer CloseConnection(conn)

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	listenForMessages(conn)
}

func listenForMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Println("Unexpected WebSocket closure:", err)
			} else {
				log.Println("Client disconnected:", err)
			}

			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		var receivedMsg Event
		err = json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		log.Println("\rServer Received:", receivedMsg)

		// Send the decoded struct to the broadcast channel
		broadcast <- receivedMsg
	}
}

func HandleBroadcast() {
	for {
		msg := <-broadcast

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("\rError marshaling JSON:", err)
			continue
		}

		mutex.Lock()
		for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, jsonMsg)
			if err != nil {
				log.Println("\rWrite error:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		mutex.Unlock()
	}
}

func BroadcastMessage(event Event) {
	message, err := json.Marshal(event)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func CloseConnection(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := clients[conn]; exists {
		log.Println("Closing WebSocket connection...")
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Goodbye!"))
		conn.Close()
		delete(clients, conn)
	}
}
