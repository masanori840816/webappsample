package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketMessage struct {
	MessageType string `json:"messageType"`
	Data        string `json:"data"`
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Close the connection when the for-loop operation is finished.
	defer conn.Close()

	message := &websocketMessage{}
	for {
		messageType, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}
		log.Println("Type: " + strconv.Itoa(messageType) + " Data: " + message.Data)
		conn.WriteJSON(message)
	}
	//client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//client.hub.register <- client
}
