package websocketcore

import (
	"encoding/json"
	"time"

	websocketevent "github.com/wang900115/LCA/internal/adapter/websocket/event"
	websocketmodel "github.com/wang900115/LCA/internal/adapter/websocket/model"
)

type Hub struct {
	Clients    map[uint]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Comment    chan websocketmodel.WSMessage
	Edit       chan websocketmodel.WSMessage
	Delete     chan websocketmodel.WSMessage
	System     chan websocketmodel.SYSMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uint]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Comment:    make(chan websocketmodel.WSMessage),
		Edit:       make(chan websocketmodel.WSMessage),
		Delete:     make(chan websocketmodel.WSMessage),
		System:     make(chan websocketmodel.SYSMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {

		// Register a new client
		case client := <-h.Register:
			if h.Clients[client.Channel.ID] == nil {
				h.Clients[client.Channel.ID] = make(map[*Client]bool)
			}
			h.Clients[client.Channel.ID][client] = true

			payload := websocketmodel.WSMessage{
				Type:      websocketevent.EVENT_USER_JOIN,
				Sender:    client.User.Username,
				Content:   "join the channel",
				Timestamp: time.Now().UTC().Unix(),
			}

			msg, _ := json.Marshal(payload)

			for c := range h.Clients[client.Channel.ID] {
				c.Send <- msg
			}

		// Unregister a client
		case client := <-h.Unregister:
			if clients, ok := h.Clients[client.Channel.ID]; ok {
				if _, ok := clients[client]; ok {

					payload := websocketmodel.WSMessage{
						Type:      websocketevent.EVENT_USER_LEAVE,
						Sender:    client.User.Username,
						Content:   "leave the channel",
						Timestamp: time.Now().UTC().Unix(),
					}

					msg, _ := json.Marshal(payload)
					for c := range h.Clients[client.Channel.ID] {
						c.Send <- msg
					}

					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.Clients, client.Channel.ID)
					}
				}
			}

		// Comment message
		case message := <-h.Comment:
			if clients, ok := h.Clients[message.ChannelID]; ok {
				event := websocketmodel.WSMessage{
					Type:      websocketevent.EVENT_USER_COMMENT,
					MessageID: message.MessageID,
					UserID:    message.UserID,
					Content:   message.Content,
					Timestamp: time.Now().UTC().Unix(),
				}
				data, _ := json.Marshal(event)
				for c := range clients {
					c.Send <- data
				}
			}

		// Edit message
		case message := <-h.Edit:
			if clients, ok := h.Clients[message.ChannelID]; ok {
				event := websocketmodel.WSMessage{
					Type:      websocketevent.EVENT_USER_EDIT,
					MessageID: message.MessageID,
					UserID:    message.UserID,
					Content:   message.Content,
					Timestamp: time.Now().UTC().Unix(),
				}
				data, _ := json.Marshal(event)
				for c := range clients {
					c.Send <- data
				}
			}

		// Delete message
		case message := <-h.Delete:
			if clients, ok := h.Clients[message.ChannelID]; ok {
				event := websocketmodel.WSMessage{
					Type:      websocketevent.EVENT_USER_DELETE,
					MessageID: message.MessageID,
					UserID:    message.UserID,
					Timestamp: time.Now().UTC().Unix(),
				}
				data, _ := json.Marshal(event)
				for c := range clients {
					c.Send <- data
				}
			}

		// System message
		case sysmsg := <-h.System:

			switch sysmsg.Type {
			// system fixed
			case websocketevent.EVENT_SYSTEM_FIX:
				event := websocketmodel.SYSMessage{
					MessageID: sysmsg.MessageID,
					Type:      websocketevent.EVENT_SYSTEM_GLOBAL,
					Message:   sysmsg.Message,
					Timestamp: time.Now().UTC().Unix(),
				}
				data, _ := json.Marshal(event)
				for _, clients := range h.Clients {
					for c := range clients {
						c.Send <- data
					}
				}
			// channel deleted
			case websocketevent.EVENT_SYSTEM_CHANNEL_DELETE:
				if clients, ok := h.Clients[sysmsg.ChannelID]; ok {
					event := websocketmodel.SYSMessage{
						MessageID: sysmsg.MessageID,
						Type:      websocketevent.EVENT_SYSTEM_LOCAL,
						Message:   sysmsg.Message,
						Timestamp: time.Now().UTC().Unix(),
					}
					data, _ := json.Marshal(event)
					for c := range clients {
						c.Send <- data
					}
				}
			//	channel fixed
			case websocketevent.EVENT_SYSTEM_CHANNEL_FIX:
				if clients, ok := h.Clients[sysmsg.ChannelID]; ok {
					event := websocketmodel.SYSMessage{
						MessageID: sysmsg.MessageID,
						Type:      websocketevent.EVENT_SYSTEM_LOCAL,
						Message:   sysmsg.Message,
						Timestamp: time.Now().UTC().Unix(),
					}
					data, _ := json.Marshal(event)
					for c := range clients {
						c.Send <- data
					}
				}
			// global event
			case websocketevent.EVENT_SYSTEM_GLOBAL:
				event := websocketmodel.SYSMessage{
					MessageID: sysmsg.MessageID,
					Type:      websocketevent.EVENT_SYSTEM_GLOBAL,
					Message:   sysmsg.Message,
					Timestamp: time.Now().UTC().Unix(),
				}
				data, _ := json.Marshal(event)
				for _, clients := range h.Clients {
					for c := range clients {
						c.Send <- data
					}
				}
			//	local event(specify channel)
			case websocketevent.EVENT_SYSTEM_LOCAL:
				if clients, ok := h.Clients[sysmsg.ChannelID]; ok {
					event := websocketmodel.SYSMessage{
						MessageID: sysmsg.MessageID,
						Type:      websocketevent.EVENT_SYSTEM_LOCAL,
						Message:   sysmsg.Message,
						Timestamp: time.Now().UTC().Unix(),
					}
					data, _ := json.Marshal(event)
					for c := range clients {
						c.Send <- data
					}
				}
			}
		}
	}
}
