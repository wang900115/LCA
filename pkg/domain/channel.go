package domain

import (
	"time"
)

type Channel struct {
	ID          uint   `json:"id"`
	ChannelName string `json:"channelName"`
	ChannelType string `json:"channelType"`
	Users       []User `json:"users"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
