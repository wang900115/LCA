package model

import (
	"LCA/internal/domain/entities"
	"time"
)

type Message struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UUID        string `json:"uuid" gorm:"unique;not null"`
	ChannelUUID string `json:"channel_uuid" gorm:"not null"`
	UserUUID    string `json:"user_uuid" gorm:"not null"`
	Content     string `json:"content" gorm:"not null"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (m Message) TableName() string {
	return "messages"
}

func (m Message) ToDomain() entities.Message {
	return entities.Message{
		UUID:        m.UUID,
		ChannelUUID: m.ChannelUUID,
		UserUUID:    m.UserUUID,
		Content:     m.Content,
	}
}
