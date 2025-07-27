package domain

import "time"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password *string

	Role        string `json:"role"`
	FirstEmail  string `json:"first_email"`
	SecondEmail string `json:"second_email"`

	Phone     string    `json:"phone"`
	NickName  string    `json:"nickName"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Birth     time.Time `json:"birth"`
	Country   string    `json:"country"`
	City      string    `json:"city"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
