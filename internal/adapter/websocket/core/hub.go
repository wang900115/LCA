package websocketcore

import (
	"encoding/json"
	"time"

	websocketevent "github.com/wang900115/LCA/internal/adapter/websocket/event"
	websocketmodel "github.com/wang900115/LCA/internal/adapter/websocket/model"
)

type BroadcastMessage struct {
	ChannelID uint
	Message   []byte
}

type Hub struct {
	Clients    map[uint]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint]map[*Client]bool),
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
			if h.Clients[client.ChannelID] == nil {
				h.Clients[client.ChannelID] = make(map[*Client]bool)
			}
			h.Clients[client.ChannelID][client] = true

			payload := websocketmodel.ChannelPayload{
				Type:      websocketevent.EVENT_USER_JOIN,
				Sender:    client.User.Username,
				Content:   "join the channel",
				Timestamp: time.Now().UTC().Unix(),
			}

			msg, _ := json.Marshal(payload)

			for c := range h.Clients[client.ChannelID] {
				c.Send <- msg
			}

		// Unregister a client
		case client := <-h.Unregister:
			if clients, ok := h.Clients[client.ChannelID]; ok {
				if _, ok := clients[client]; ok {

					payload := websocketmodel.ChannelPayload{
						Type:      websocketevent.EVENT_USER_LEAVE,
						Sender:    client.User.Username,
						Content:   "leave the channel",
						Timestamp: time.Now().UTC().Unix(),
					}

					msg, _ := json.Marshal(payload)
					for c := range h.Clients[client.ChannelID] {
						c.Send <- msg
					}

					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.Clients, client.ChannelID)
					}
				}
			}

		// Broadcast a message to all clients in the channel
		case message := <-h.Broadcast:
			if clients, ok := h.Clients[message.ChannelID]; ok {
				for c := range clients {
					c.Send <- message.Message
				}
			}
		}
	}
}
