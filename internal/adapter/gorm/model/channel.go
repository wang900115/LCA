package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

const (
	CHANNEL_TYPE_PUBLIC  = "public"
	CHANNEL_TYPE_PRIVATE = "private"
)

type Channel struct {
	gorm.Model
	Name        string    `gorm:"column:name;unique;not null"`
	FounderID   uint      `gorm:"column:founder_id; not null"`
	Founder     User      `gorm:"foreignKey:FounderID"`
	ChannelType string    `gorm:"column:channel_type;type:enum('public','private');not null"`
	Messages    []Message `gorm:"foreignKey:ChannelID"`
	Users       []User    `gorm:"foreignKey:UserID"`
}

func (c Channel) TableName() string {
	return "channel"
}

func (c Channel) ToDomain() *entities.Channel {
	var messages []*entities.Message
	for _, message := range c.Messages {
		messages = append(messages, message.ToDomain())
	}
	var users []*entities.User
	for _, user := range c.Users {
		users = append(users, user.ToDomain())
	}
	return &entities.Channel{
		Name:        c.Name,
		FounderID:   c.FounderID,
		Founder:     c.Founder.ToDomain(),
		ChannelType: c.ChannelType,
		Messages:    messages,
		Users:       users,
	}
}
