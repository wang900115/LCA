package entities

type Message struct {
	UUID        string `json:"uuid"`
	ChannelUUID string `json:"channel_uuid"`
	UserUUID    string `json:"user_uuid"`
	Content     string `json:"content"`
}
