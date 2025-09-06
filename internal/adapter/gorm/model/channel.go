package gormmodel

import (
	"time"

	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	Name        string `gorm:"column:name;unique;not nill"`
	Founder     string `gorm:"column:founder"`
	ChannelType string `gorm:"column:channel_type"`
}

func (c Channel) TableName() string {
	return "channel"
}

func (c Channel) ToDomain() *entities.Channel {
	return &entities.Channel{
		Name:        c.Name,
		Founder:     c.Founder,
		ChannelType: c.ChannelType,
	}
}

type ChannelUser struct {
	ChannelID uint    `gorm:"column:channel_id"`
	Users     []*User `gorm:"foreignKey:ChannelID;references:ChannelID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (cu ChannelUser) TableName() string {
	return "channel_user"
}

func (cu ChannelUser) ToDomain() *entities.ChannelUser {
	usersDomain := make([]*entities.User, 0, len(cu.Users))
	for _, user := range cu.Users {
		usersDomain = append(usersDomain, user.ToDomain())
	}
	return &entities.ChannelUser{
		ChannelID: cu.ChannelID,
		Users:     usersDomain,
	}
}

type ChannelMessage struct {
	ChannelID uint       `gorm:"column:channel_id"`
	Messages  []*Message `gorm:"foreignKey:ChannelID;references:ChannelID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (cm ChannelMessage) TableName() string {
	return "channel_message"
}

func (cm ChannelMessage) ToDomain() *entities.ChannelMessage {
	messagesDomain := make([]*entities.Message, 0, len(cm.Messages))
	for _, message := range cm.Messages {
		messagesDomain = append(messagesDomain, message.ToDomain())
	}
	return &entities.ChannelMessage{
		ChannelID: cm.ChannelID,
		Message:   messagesDomain,
	}
}
