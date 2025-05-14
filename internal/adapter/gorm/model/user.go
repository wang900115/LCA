package model

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID        string `json:"uuid" gorm:"unique;not null;index"`
	Username    string `json:"username" gorm:"not null;unique"`
	ChannelUUID string `json:"channel_uuid" gorm:"not null;index"`

	Channel Channel `gorm:"foreignkey:ChannelUUID; references:UUID"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) ToDomain() entities.User {
	return entities.User{
		UUID:        u.UUID,
		Username:    u.Username,
		ChannelUUID: u.ChannelUUID,
	}
}
