package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, you should verify the origin of the request:
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (for testing purposes)
	},
}

// wsHandler upgrades the HTTP connection to a WebSocket and echoes messages.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected", r.Header.Get("Authorization"))
	//parse users id from jwt and get users conversations from db
	//does conversations exists on redis if its not create the key with this userId if yes add this user Id to conversation
	// set the user Id to online user list  with ttl 15 seconds

	// Continuously read messages from the client.
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		//ping mesajı geldiğinde user id online listte set et client 15 saniyede bir ping atmalı
		log.Printf("Received: %s", message)

		// Echo the message back to the client.
		if err = conn.WriteMessage(messageType, message); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
