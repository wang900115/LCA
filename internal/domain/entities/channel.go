package entities

type Channel struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	FounderID   uint       `json:"founderID"`
	Founder     *User      `json:"founder"`
	ChannelType string     `json:"type"`
	Messages    []*Message `json:"messages"`
	Users       []*User    `json:"users"`
}
