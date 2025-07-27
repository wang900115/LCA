package validator

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required, min=4, max=12, alphanum"`
	Password string `json:"password" binding:"required, min=6, max=12"`

	FirstEmail  string `json:"firstEmail" binding:"required, email"`
	SecondEmail string `json:"secondEmail" binding:"omitempty, email"`
	Phone       string `json:"phone" binding:"required,e164"`
	NickName    string `json:"nickName" binding:"required,  max=10"`
	FirstName   string `json:"firstName" binding:"required, max=30"`
	LastName    string `json:"lastName" binding:"required, max=30"`

	Birth   time.Time `json:"birth" binding:"required"`
	Country string    `json:"country" binding:"required"`
	City    string    `json:"city" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty, min=4, max=12, alphanum"`
	Password string `json:"password" binding:"omitempty, min=6, max=12"`

	FirstEmail  string `json:"firstEmail" binding:"required, email"`
	SecondEmail string `json:"secondEmail" binding:"omitempty, email"`
	Phone       string `json:"phone" binding:"required,e164"`
	NickName    string `json:"nickName" binding:"omitempty,  max=10"`
	FirstName   string `json:"firstName" binding:"omitempty, max=30"`
	LastName    string `json:"lastName" binding:"omitempty, max=30"`

	Birth   time.Time `json:"birth" binding:"omitempty"`
	Country string    `json:"country" binding:"omitempty"`
	City    string    `json:"city" binding:"omitempty"`
}

type DeleteUserRequest struct {
	Confirm string `json:"confirm" binding:"required,confirm"`
}

type QueryUserRequest struct {
}

type LoginAuthRequest struct {
	Username string `json:"username" binding:"omitempty, min=4, max=30, alphanum"`
	Password string `json:"password" binding:"required, min=6, max=20"`
}

type LogoutAuthRequest struct {
}

func Confirm(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return strings.EqualFold(value, "I AM SURE")
}
