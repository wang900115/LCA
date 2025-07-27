package gormmodel

import (
	"time"

	"github.com/wang900115/LCA/pkg/domain"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(50);not null;unique;comment:用戶帳號"`
	Password string `gorm:"column:password;type:char(60);not null;comment:用戶密碼希哈"`

	Role        string    `gorm:"column:role;type:varchar(50);default:'general';comment:權限"`
	FirstEmail  string    `gorm:"column:first_email;type:varchar(100);not null;unique; comment:電子郵件"`
	SecondEmail string    `gorm:"column:second_email;type:varchar(100);comment:備份電子郵件"`
	Phone       string    `gorm:"column:phone;type:varchar(100); not null;unique; comment:電話"`
	NickName    string    `gorm:"column:nick_name;type:varchar(50);comment:暱稱"`
	FirstName   string    `gorm:"column:first_name;type:varchar(50);comment:首名"`
	LastName    string    `gorm:"column:last_name;type:varchar(50);comment:尾名"`
	Birth       time.Time `gorm:"column:birth;type:date;comment:出生"`
	Country     string    `gorm:"column:country;type:varchar(50);comment:國家"`
	City        string    `gorm:"column:city;type:varchar(50);comment:城市"`

	ChannelID uint    `gorm:"column:channel_id;index;comment:目前所屬頻道ID"`
	Channel   Channel `gorm:"foreignkey:ChannelID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Messages []Message `gorm:"foreignKey:UserID"`
}

/*
FOREIGN KEY (channel_id) REFERENCES channels(id)
ON UPDATE CASADE 	=> 如果 channels.id 被更新 則 users.channel_id也會跟著更新
ON DELETE SET NULL  => 如果 channels.id 被刪除 則 users.channel_id會變成 NULL
*/

func (u *User) TableName() string {
	return "user"
}

func (u User) ToDomain() domain.User {
	return domain.User{
		ID:       u.Model.ID,
		Username: u.Username,
		Password: nil,

		Role:      u.Role,
		NickName:  u.NickName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Birth:     u.Birth,
		Country:   u.Country,
		City:      u.City,

		CreatedAt: u.Model.CreatedAt,
		UpdatedAt: u.Model.UpdatedAt,
	}
}

func (u User) FromDomain(user domain.User) User {
	return User{
		Model: gorm.Model{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},

		Username: user.Username,
		Password: *user.Password,
		Role:     user.Role,

		FirstEmail:  user.FirstEmail,
		SecondEmail: user.SecondEmail,
		Phone:       user.Phone,

		NickName:  user.NickName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Birth:     user.Birth,
		Country:   user.Country,
		City:      user.City,
	}
}
