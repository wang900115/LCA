package model

import (
	"LCA/internal/domain/entities"
	"time"
)

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UUID        string `json:"uuid" gorm:"unique;not null"`
	Username    string `json:"username" gorm:"not null"`
	ChannelUUID string `json:"channel_uuid" gorm:"not null;index"`
	Status      string `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
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
