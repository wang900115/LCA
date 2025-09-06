package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Sender  uint   `gorm:"column:sender;not null"`
	Content string `gorm:"column:content;not null"`
}

func (m Message) TableName() string {
	return "message"
}

func (m Message) ToDomain() *entities.Message {
	return &entities.Message{
		ID:      m.ID,
		Content: m.Content,
	}
}
