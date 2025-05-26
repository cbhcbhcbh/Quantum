package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/ws/connect", nil)
	if err != nil {
		panic("Failed to connect to WebSocket server: " + err.Error())
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	if err != nil {
		panic("Failed to send message: " + err.Error())
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		panic("Failed to read message: " + err.Error())
	}
	log.Println("Received message:", string(p))
}
