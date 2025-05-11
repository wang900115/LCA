package connection

import (
	"LCA/internal/adapter/websocket/event"
	"encoding/json"
	"time"
)

type BroadcastMessage struct {
	ChannelUUID string
	Message     []byte
}

type Hub struct {
	Clients    map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {

		// Register a new client
		case client := <-h.Register:
			if h.Clients[client.ChannelUUID] == nil {
				h.Clients[client.ChannelUUID] = make(map[*Client]bool)
			}
			h.Clients[client.ChannelUUID][client] = true

			payload := event.MessagePayload{
				Type:      event.EventJoin,
				Sender:    client.Username,
				Content:   "join the channel",
				Timestamp: time.Now().Format(time.RFC3339),
			}
			msg, _ := json.Marshal(payload)

			for c := range h.Clients[client.ChannelUUID] {
				c.Send <- msg
			}

		// Unregister a client
		case client := <-h.Unregister:
			if clients, ok := h.Clients[client.ChannelUUID]; ok {
				if _, ok := clients[client]; ok {

					payload := event.MessagePayload{
						Type:      event.EventLeave,
						Sender:    client.Username,
						Content:   "leave the channel",
						Timestamp: time.Now().Format(time.RFC3339),
					}

					msg, _ := json.Marshal(payload)
					for c := range h.Clients[client.ChannelUUID] {
						c.Send <- msg
					}

					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.Clients, client.ChannelUUID)
					}
				}
			}
		// Broadcast a message to all clients in the channel
		case message := <-h.Broadcast:
			if clients, ok := h.Clients[message.ChannelUUID]; ok {
				for c := range clients {
					c.Send <- message.Message
				}
			}
		}
	}
}
