package model

import (
	"LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UUID        string `json:"uuid" gorm:"unique;not null;index"`
	ChannelUUID string `json:"channel_uuid" gorm:"not null"`
	UserUUID    string `json:"user_uuid" gorm:"not null"`
	Content     string `json:"content" gorm:"not null"`

	Channel Channel `gorm:"foreignKey:ChannelUUID;references:UUID"`
	User    User    `gorm:"foreignKey:UserUUID;references:UUID"`
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
