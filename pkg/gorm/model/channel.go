package gormmodel

import (
	"github.com/wang900115/LCA/pkg/domain"
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	ChannelName string `gorm:"column:channel_name;tpye:varchar(50);unique;not null;comment:頻道名稱"`
	ChannelType string `gorm:"column:channel_type;tpye:varchar(30);not null;comment頻道類別"`

	Users []User `gorm:"foreignKey:ChannelID"`
}

func (c *Channel) TableName() string {
	return "channel"
}

func (c Channel) ToDomain() domain.Channel {

	var users []domain.User
	for _, user := range c.Users {
		users = append(users, user.ToDomain())
	}

	return domain.Channel{
		ID:          c.ID,
		ChannelName: c.ChannelName,
		ChannelType: c.ChannelType,
		Users:       users,

		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c Channel) FromDomain(channel domain.Channel) Channel {

	var users []User

	for _, user := range channel.Users {
		users = append(users, User{}.FromDomain(user))
	}

	return Channel{
		Model: gorm.Model{
			ID:        channel.ID,
			CreatedAt: channel.CreatedAt,
			UpdatedAt: channel.UpdatedAt,
		},
		Users:       users,
		ChannelName: channel.ChannelName,
		ChannelType: channel.ChannelType,
	}
}
