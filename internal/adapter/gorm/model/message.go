package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ChannelID uint   `gorm:"index; not null"`
	Sender    uint   `gorm:"column:sender;not null"`
	Content   string `gorm:"column:content;not null"`

	Channel Channel `gorm:"foreignKey:ChannelID"`
	User    User    `gorm:"foreignKey:Sender"`
}

func (m Message) TableName() string {
	return "message"
}

func (m Message) ToDomain() *entities.Message {
	return &entities.Message{
		ID:      m.ID,
		Sender:  m.Sender,
		Content: m.Content,
	}
}
