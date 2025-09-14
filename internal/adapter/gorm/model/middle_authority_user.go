package gormmodel

import (
	"gorm.io/gorm"
)

type MiddleAuthorityUser struct {
	gorm.Model
	UserID      uint `gorm:"primaryKey"`
	AuthorityID uint `gorm:"primaryKey"`

	Authority Authority `gorm:"foreignKey:AuthorityID;reference:ID"`
	User      User      `gorm:"foreignKey:UserID;reference:ID"`
}

func (mau MiddleAuthorityUser) TableName() string {
	return "user_authority_middle"
}
