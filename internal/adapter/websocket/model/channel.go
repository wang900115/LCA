package websocketmodel

type WSMessage struct {
	MessageID uint
	ChannelID uint
	UserID    uint
	Type      string
	Sender    string
	Content   string
	Timestamp int64
}

type SYSMessage struct {
	MessageID uint
	ChannelID uint
	Type      string
	Message   []byte
	Timestamp int64
}
