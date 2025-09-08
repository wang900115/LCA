package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"gorm.io/gorm"
)

type MiddleChannelUser struct {
	gorm.Model
	UserID    uint   `gorm:"primaryKey"`
	ChannelID uint   `gorm:"primaryKey"`
	Role      string `gorm:"column:role"`
	LastJoin  int64  `gorm:"column:last_join"`

	Channel Channel `gorm:"foreignKey:ChannelID;reference:ID"`
	User    User    `gorm:"foreignKey:UserID;reference:ID"`
}

func (mcu MiddleChannelUser) TableName() string {
	return "user_channel_middle"
}

func (mcu *MiddleChannelUser) ToDomain() *entities.UserJoin {
	return &entities.UserJoin{
		Role:     mcu.Role,
		LastJoin: mcu.LastJoin,
	}
}
