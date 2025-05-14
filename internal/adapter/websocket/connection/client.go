package connection

import (
	"encoding/json"
	"time"

	"github.com/wang900115/LCA/internal/adapter/websocket/event"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserUUID    string
	ChannelUUID string
	Username    string
	Conn        *websocket.Conn
	Send        chan []byte

	messageUsecase *usecase.MessageUsecase
}

func NewClient(userUUID, channelUUID, userName string, conn *websocket.Conn, messageUsecase *usecase.MessageUsecase) *Client {
	return &Client{
		UserUUID:       userUUID,
		ChannelUUID:    channelUUID,
		Username:       userName,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		messageUsecase: messageUsecase,
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
		var payload event.MessagePayload
		if err := json.Unmarshal(message, &payload); err != nil {
			continue
		}

		payload.Type = event.EventMessage
		payload.Sender = c.Username
		payload.Timestamp = time.Now().Format(time.RFC3339)

		messageUUID, err := c.messageUsecase.CreateMessage(c.ChannelUUID, c.UserUUID, payload.Content)

		if err != nil {
			continue
		}

		encoded, err := json.Marshal(payload)
		if err != nil {
			continue
		}

		hub.Broadcast <- BroadcastMessage{
			ChannelUUID: c.ChannelUUID,
			Message:     encoded,
		}

		payload.UUID = messageUUID
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
