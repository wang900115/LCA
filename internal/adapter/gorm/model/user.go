package gormmodel

import (
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"column:user_name;not null;unique"`
	Password string `gorm:"column:password;not null;unique"`
	NickName string `gorm:"column:nick_name;not null"`
	FullName string `gorm:"column:full_name"`
	LastName string `gorm:"column:last_name"`
	Email    string `gorm:"column:email"`
	Phone    string `gorm:"column:phone"`

	Birth  int64  `gorm:"column:birth;not null"`
	Status string `gorm:"column:status;default:'active'"`
}

func (u User) TableName() string {
	return "user"
}

func (u User) ToDomain() *entities.User {
	return &entities.User{
		ID:       u.ID,
		Username: u.Username,
		Password: nil,
		NickName: u.NickName,
		FullName: u.FullName,
		LastName: u.LastName,
		Email:    u.Email,
		Phone:    u.Phone,
		Birth:    u.Birth,
		Status:   nil,
	}
}

type UserLogin struct {
	gorm.Model
	UserID     uint   `gorm:"column:user_id"`
	LastLogin  int64  `gorm:"column:last_login"`
	IPAddress  string `gorm:"column:ip_address"`
	DeviceInfo string `gorm:"column:device_info"`
	User       User   `gorm:"foreignKey:UserID;reference:ID"`
}

func (ul UserLogin) TableName() string {
	return "user_login"
}

func (ul UserLogin) ToDomain() *entities.UserLogin {
	return &entities.UserLogin{
		LastLogin:  ul.LastLogin,
		IPAddress:  nil,
		DeviceInfo: nil,
	}
}
