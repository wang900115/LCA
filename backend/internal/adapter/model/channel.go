package model

import "time"

type Channel struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	UUID  string `json:"uuid" gorm:"unique;not null"`
	Users []User `json:"users" gorm:"foreignKey:channel_uuid;reference:uuid"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
}
