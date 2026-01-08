package main

import "log"

type ChatGroups struct {
	rooms      map[string]*SSEHub
	register   chan *SSEHub
	unregister chan *SSEHub
	close      chan int8
}

func NewChatGroups() *ChatGroups {
	return &ChatGroups{
		rooms:      make(map[string]*SSEHub),
		register:   make(chan *SSEHub),
		unregister: make(chan *SSEHub),
		close:      make(chan int8),
	}
}
func (groups *ChatGroups) Run() {
	defer func() {
		groups.Close()
	}()
	for {
		select {
		case hub := <-groups.register:
			groups.rooms[hub.name] = hub
		case hub := <-groups.unregister:
			delete(groups.rooms, hub.name)
		case i := <-groups.close:
			log.Println(i)
			return
		}
	}
}
func (groups *ChatGroups) GetOrCreateRoom(name string) *SSEHub {
	if _, ok := groups.rooms[name]; ok {
		return groups.rooms[name]
	}
	newRoom := NewSSEHub(name)
	groups.register <- newRoom
	go newRoom.run(groups.unregister)
	return newRoom
}
func (groups *ChatGroups) Close() {
	for _, r := range groups.rooms {
		r.close()
	}
	close(groups.register)
	close(groups.unregister)
}
