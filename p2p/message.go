package p2p

import "encoding/json"

type Message struct {
	ID        string      // Unique identifier for the message
	ChannelID string      // ID of the channel the message belongs to
	Receiver  string      // ID of the message receiver
	Type      string      // Type of message (e.g., "text", "file", etc.)
	Payload   interface{} // Actual message content or data
}

func NewMessage(id, channelID, receiver, msgType string, payload interface{}) *Message {
	return &Message{
		ID:        id,
		ChannelID: channelID,
		Receiver:  receiver,
		Type:      msgType,
		Payload:   payload,
	}
}

func (m *Message) GetID() string {
	return m.ID
}

func (m *Message) GetChannelID() string {
	return m.ChannelID
}

func (m *Message) GetReceiver() string {
	return m.Receiver
}

func (m *Message) GetType() string {
	return m.Type
}

func (m *Message) GetPayload() interface{} {
	return m.Payload
}

func (m *Message) CallbackInternalService() {
	// Placeholder for future implementation
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}
