package connection

type BroadcastMessage struct {
	ChannelUUID string
	Message     []byte
	Sender      *Client
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
		case client := <-h.Register:
			if h.Clients[client.ChannelUUID] == nil {
				h.Clients[client.ChannelUUID] = make(map[*Client]bool)
			}
			h.Clients[client.ChannelUUID][client] = true

		case client := <-h.Unregister:
			if clients, ok := h.Clients[client.ChannelUUID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Clients, client.ChannelUUID)
					}
				}
			}
		case message := <-h.Broadcast:
			if clients, ok := h.Clients[message.ChannelUUID]; ok {
				for c := range clients {
					c.Send <- message.Message
				}
			}
		}
	}
}
