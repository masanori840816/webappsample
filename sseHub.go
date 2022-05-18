package main

import "encoding/json"

type SSEHub struct {
	clients    map[*SSEClient]bool
	broadcast  chan ClientMessage
	register   chan *SSEClient
	unregister chan *SSEClient
}

func newSSEHub() *SSEHub {
	return &SSEHub{
		clients:    make(map[*SSEClient]bool),
		broadcast:  make(chan ClientMessage),
		register:   make(chan *SSEClient),
		unregister: make(chan *SSEClient),
	}
}
func (h *SSEHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.messageChan)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			m, _ := json.Marshal(message)
			jsonText := string(m)

			for client := range h.clients {
				if client.userName == message.UserName {
					continue
				}
				select {
				case client.messageChan <- jsonText:
				default:
					close(client.messageChan)
					delete(h.clients, client)
				}
			}
		}
	}
}
