package model

import (
	"LCA/internal/domain/entities"
	"time"
)

type Channel struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	UUID  string `json:"uuid" gorm:"unique;not null"`
	Users []User `json:"users" gorm:"foreignKey:channel_uuid;reference:uuid"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c Channel) TableName() string {
	return "channels"
}

func (c Channel) ToDomain() entities.Channel {
	usersDomain := make([]entities.User, len(c.Users))
	for _, user := range c.Users {
		usersDomain = append(usersDomain, user.ToDomain())
	}
	return entities.Channel{
		UUID:  c.UUID,
		Users: usersDomain,
	}
}
