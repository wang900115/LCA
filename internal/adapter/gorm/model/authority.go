package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"gorm.io/gorm"
)

type Authority struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (Authority) TableName() string {
	return "authority"
}

func (a Authority) ToDomain() *entities.Authority {
	return &entities.Authority{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
	}
}
