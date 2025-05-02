package model

import "time"

type Message struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UUID        string `json:"uuid" gorm:"unique;not null"`
	ChannelUUID string `json:"channel_uuid" gorm:"not null"`
	Sender      string `json:"sender" gorm:"not null"`
	Content     string `json:"content" gorm:"not null"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
}
