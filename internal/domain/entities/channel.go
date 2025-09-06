package entities

type Channel struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Founder     string `json:"founder"`
	ChannelType string `json:"type"`
}

type ChannelUser struct {
	ChannelID uint    `json:"channel_id"`
	Users     []*User `json:"users"`
}

type ChannelMessage struct {
	ChannelID uint       `json:"channel_id"`
	Message   []*Message `json:"messages"`
}
