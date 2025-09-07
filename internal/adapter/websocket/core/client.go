package websocketcore

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	ChannelID uint
	Send      chan []byte
}

func NewClient(conn *websocket.Conn, channelId uint) *Client {
	return &Client{
		Conn:      conn,
		ChannelID: channelId,
		Send:      make(chan []byte, 256),
	}
}

// Reads messages from websocket connection and sends them to the hub's broadcast channel.
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// hub.Broadcast <- BroadcastMessage{
		// 	hub.Clients
		// }
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
