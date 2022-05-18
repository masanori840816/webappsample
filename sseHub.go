package main

import "encoding/json"

type SseHub struct {
	clients    map[*SseClient]bool
	broadcast  chan ClientMessage
	register   chan *SseClient
	unregister chan *SseClient
}

func newSseHub() *SseHub {
	return &SseHub{
		clients:    make(map[*SseClient]bool),
		broadcast:  make(chan ClientMessage),
		register:   make(chan *SseClient),
		unregister: make(chan *SseClient),
	}
}
func (h *SseHub) run() {
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
