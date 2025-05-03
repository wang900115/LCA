package connection

import "github.com/gorilla/websocket"

type Client struct {
	UserUUID    string
	ChannelUUID string
	Username    string
	Conn        *websocket.Conn
	Send        chan []byte
}

func NewClient(userUUID, channelUUID, userName string, conn *websocket.Conn) *Client {
	return &Client{
		UserUUID:    userUUID,
		ChannelUUID: channelUUID,
		Username:    userName,
		Conn:        conn,
		Send:        make(chan []byte, 256),
	}
}

// Reads messages from websocket connection and sends them to the hub's broadcast channel.
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		hub.Broadcast <- BroadcastMessage{
			ChannelUUID: c.ChannelUUID,
			Message:     message,
			Sender:      c,
		}
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
