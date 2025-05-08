package model

import (
	"LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID        string `json:"uuid" gorm:"unique;not null"`
	Username    string `json:"username" gorm:"not null"`
	ChannelUUID string `json:"channel_uuid" gorm:"not null;index"`
	Status      string `json:"status"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) ToDomain() entities.User {
	return entities.User{
		UUID:        u.UUID,
		Username:    u.Username,
		ChannelUUID: u.ChannelUUID,
		Status:      u.Status,
	}
}
