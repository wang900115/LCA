package websocketcore

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	websocketevent "github.com/wang900115/LCA/internal/adapter/websocket/event"
	websocketmodel "github.com/wang900115/LCA/internal/adapter/websocket/model"
	"github.com/wang900115/LCA/internal/domain/entities"
)

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Channel entities.Channel
	User    entities.User
	Send    chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn, channel entities.Channel, user entities.User) *Client {
	return &Client{
		Hub:     hub,
		Conn:    conn,
		Channel: channel,
		User:    user,
		Send:    make(chan []byte, 256),
	}
}

func (c *Client) HandleMessage(msg websocketmodel.WSMessage) {
	switch msg.Type {
	case websocketevent.EVENT_USER_COMMENT:
		c.Hub.Comment <- msg
	case websocketevent.EVENT_USER_EDIT:
		c.Hub.Edit <- msg
	case websocketevent.EVENT_USER_DELETE:
		c.Hub.Delete <- msg
	}
}

// Reads messages from websocket connection and sends them to the hub's responding channel.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg websocketmodel.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			// todo
			continue
		}
		c.HandleMessage(msg)
	}
}

// Writes messages to the frontend
func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
