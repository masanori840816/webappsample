package main

type SseHub struct {
	// registered clients
	clients    map[*SseClient]bool
	broadcast  chan string
	register   chan *SseClient
	unregister chan *SseClient
}

func newSseHub() *SseHub {
	return &SseHub{
		clients:    make(map[*SseClient]bool),
		broadcast:  make(chan string),
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
			for client := range h.clients {
				select {
				case client.messageChan <- message:
				default:
					close(client.messageChan)
					delete(h.clients, client)
				}
			}
		}
	}
}
