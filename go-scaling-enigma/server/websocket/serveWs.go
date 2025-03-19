package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan string)
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
		log.Println("\rServer Received:", string(msg))

		broadcast <- string(msg)
	}

}

func HandleBroadcast() {
	for {
		msg := <-broadcast

		mutex.Lock()
		for conn := range clients {
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("\rWrite error:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		mutex.Unlock()
	}
}

func BroadcastMessage(message string) {
	broadcast <- message
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
