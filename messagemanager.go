package main

import "log"

type MessageManager struct {
	clients   map[string]*Client
	Broadcast chan []byte
	Login     chan *Client
	Logout    chan *Client
}

func NewMessageManager() *MessageManager {
	return &MessageManager{
		clients:   make(map[string]*Client, 100),
		Broadcast: make(chan []byte, 100),
		Login:     make(chan *Client, 100),
		Logout:    make(chan *Client, 100),
	}
}

func (mm *MessageManager) run() {
	for {
		select {
		case client := <-mm.Login:
			log.Println(client)
			mm.clients[client.UID] = client
		case client := <-mm.Logout:
			delete(mm.clients, client.UID)
		case message := <-mm.Broadcast:
			log.Println("broadcast: ", string(message))
			for uid := range mm.clients {
				client := mm.clients[uid]
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(mm.clients, uid)
				}
			}
		}
	}
}
