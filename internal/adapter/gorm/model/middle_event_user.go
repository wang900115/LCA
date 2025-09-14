package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"gorm.io/gorm"
)

type MiddleEventUser struct {
	gorm.Model
	UserID        uint   `gorm:"primaryKey"`
	EventID       uint   `gorm:"primaryKey"`
	Role          string `gorm:"column:role"`
	LastParticate int64  `gorm:"column:last_join"`

	Event Event `gorm:"foreignKey:EventID;reference:ID"`
	User  User  `gorm:"foreignKey:UserID;reference:ID"`
}

func (meu MiddleEventUser) TableName() string {
	return "user_event_middle"
}

func (meu *MiddleEventUser) ToDomain() *entities.UserParticate {
	return &entities.UserParticate{
		Role:          meu.Role,
		LastParticate: meu.LastParticate,
	}
}
