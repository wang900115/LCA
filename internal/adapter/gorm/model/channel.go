package model

import (
	"LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	UUID  string `json:"uuid" gorm:"unique;not null;index"`
	Users []User `json:"users" gorm:"foreignKey:ChannelUUID;references:UUID"`
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
		UUID:  c.UUID,
		Users: usersDomain,
	}
}
