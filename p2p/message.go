package p2p

type Message struct {
	ID        string      // Unique identifier for the message
	ChannelID string      // ID of the channel the message belongs to
	Receiver  string      // ID of the message receiver
	Type      string      // Type of message (e.g., "text", "file", etc.)
	Payload   interface{} // Actual message content or data
}
