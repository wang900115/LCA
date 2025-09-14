package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"gorm.io/gorm"
)

const (
	EVENT_TYPE_GLOBAL = "global"
	EVENT_TYPE_LOCAL  = "local"
)

type Event struct {
	gorm.Model
	Name      string `gorm:"column:name;unique;not null"`
	EventType string `gorm:"column:event_type;type:enum('global','local');not null"`
	StartTime int64  `gorm:"column:start_time"`
	EndTime   int64  `gorm:"column:end_time"`
	Users     []User `gorm:"foreignKey:UserID"`
}

func (e Event) TableName() string {
	return "event"
}

func (e Event) ToDomain() *entities.Event {
	var users []*entities.User
	for _, user := range e.Users {
		users = append(users, user.ToDomain())
	}
	return &entities.Event{
		ID:        e.ID,
		Name:      e.Name,
		EventType: e.EventType,
		StartTime: e.StartTime,
		EndTime:   e.EndTime,
		Users:     users,
	}
}
