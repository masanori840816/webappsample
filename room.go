package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	listLock    sync.RWMutex
	connections []connectionState
)

type websocketMessage struct {
	MessageType string `json:"messageType"`
	Data        string `json:"data"`
}
type connectionState struct {
	websocket *threadSafeWriter
}
type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &threadSafeWriter{unsafeConn, sync.Mutex{}}
	// Close the connection when the for-loop operation is finished.
	defer conn.Close()
	listLock.Lock()
	connections = append(connections, connectionState{websocket: conn})
	listLock.Unlock()

	message := &websocketMessage{}
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}
		for _, c := range connections {
			c.websocket.WriteJSON(message)
		}
	}
}
func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}
