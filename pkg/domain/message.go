package domain

import "time"

type Message struct {
	ID uint `json:"id"`

	ChannelID uint `json:"channelID"`
	UserID    uint `json:"userID"`

	MsgType string `json:"msgType"`
	Status  string `json:"status"`

	ReplyToID *uint `json:"replyToID"`

	Content   string `json:"content"`
	AttachURL string `json:"attachURL"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
