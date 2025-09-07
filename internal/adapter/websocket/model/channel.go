package websocketmodel

type ChannelPayload struct {
	Type      string
	Sender    string
	Content   string
	Timestamp int64
}
