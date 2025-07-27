package gormmodel

import (
	"github.com/wang900115/LCA/pkg/domain"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	ChannelID uint   `gorm:"column:channel_id;not null;index;comment:所屬頻道ID"`
	UserID    uint   `gorm:"column:user_id;not null;index;comment:發送的人ID"`
	MsgType   string `gorm:"column:msg_type;varchar(20);not null;default:'text'; comment:訊息類型(text, image, viedo, file)"`
	Status    string `gorm:"column:status;type:varchar(20);not null;default:'sent'; comment:訊息狀態(sent, pending, delivered)"`

	ReplyToID *uint `gorm:"column:reply_to_id;comment:要回復的訊息ID"`

	Content   string `gorm:"column:content;type:text;not null;comment:訊息內容"`
	AttachURL string `gorm:"column:attach_url;type:varchar(255);comment:附件URL"`
}

func (m *Message) TableName() string {
	return "message"
}

func (m Message) ToDomain() domain.Message {
	return domain.Message{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		UserID:    m.UserID,

		MsgType: m.MsgType,
		Status:  m.Status,

		ReplyToID: m.ReplyToID,

		Content:   m.Content,
		AttachURL: m.AttachURL,
	}
}

func (m Message) FromDomain(message domain.Message) Message {
	return Message{
		Model: gorm.Model{
			ID:        message.ID,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		},
		ChannelID: message.ChannelID,
		UserID:    message.UserID,

		MsgType: message.MsgType,
		Status:  message.Status,

		ReplyToID: message.ReplyToID,

		Content:   m.Content,
		AttachURL: m.AttachURL,
	}
}
