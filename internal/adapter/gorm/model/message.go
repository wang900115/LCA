package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	ChannelName string `json:"channel_name" gorm:"not null"`
	Username    string `json:"user_name" gorm:"not null"`

	Content string `json:"content" gorm:"not null"`
}

func (m Message) TableName() string {
	return "messages"
}

func (m Message) ToDomain() entities.Message {
	return entities.Message{
		Channel: m.ChannelName,
		User:    m.Username,
		Content: m.Content,
	}
}
