package network

type Message struct {
	ID        string      // Unique identifier for the message
	Type      string      // Type of message (e.g., "text", "file", etc.)
	Payload   interface{} // Actual message content or data
	CreatedAt int64       // Timestamp of message creation
}

func NewMessage(id, msgType string, payload interface{}) *Message {
	return &Message{
		ID:      id,
		Type:    msgType,
		Payload: payload,
	}
}
