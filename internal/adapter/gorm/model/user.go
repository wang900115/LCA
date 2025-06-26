package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `json:"username" gorm:"not null;unique"`
	ChannelName string `json:"channel_name" gorm:"not null;index"`

	Channel  Channel   `gorm:"foreignkey:ChannelName; references:Name"`
	Messages []Message `gorm:"foreignKey:Username;references:Username"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) ToDomain() entities.User {
	return entities.User{
		Username: u.Username,
		Channel:  u.ChannelName,
	}
}
