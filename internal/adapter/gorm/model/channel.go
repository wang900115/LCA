package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	Name  string `json:"name" gorm:"unique;not nill"`
	Users []User `json:"users" gorm:"foreignKey:ChannelName;references:Name"`
}

func (c Channel) TableName() string {
	return "channels"
}

func (c Channel) ToDomain() entities.Channel {
	usersDomain := make([]entities.User, 0, len(c.Users))
	for _, user := range c.Users {
		usersDomain = append(usersDomain, user.ToDomain())
	}
	return entities.Channel{
		Name:  c.Name,
		Users: usersDomain,
	}
}
