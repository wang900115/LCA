package redismodel

import (
	"strconv"
	"time"

	"github.com/wang900115/LCA/pkg/domain"
	rediskey "github.com/wang900115/LCA/pkg/redis/key"
)

type Message struct {
	ChannelID uint
	UserID    uint

	MsgType string
	Status  string

	ReplyToID *uint

	Content   string
	AttachURL string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m Message) ToHash() map[string]interface{} {
	return map[string]interface{}{
		rediskey.REDIS_FIELD_MESSAGE_CHANNELID: m.ChannelID,
		rediskey.REDIS_FIELD_MESSAGE_USERID:    m.UserID,

		rediskey.REDIS_FIELD_MESSAGE_MSGTYPE: m.MsgType,
		rediskey.REDIS_FIELD_MESSAGE_STATUS:  m.Status,

		rediskey.REDIS_FIELD_MESSAGE_REPLYTOID: m.ReplyToID,
		rediskey.REDIS_FIELD_MESSAGE_CONTENT:   m.Content,
		rediskey.REDIS_FIELD_MESSAGE_ATTACHURL: m.AttachURL,

		rediskey.REDIS_FIELD_MESSAGE_CREATEDAT: m.CreatedAt.UTC().Format(time.RFC3339),
		rediskey.REDIS_FIELD_MESSAGE_UPDATEDAT: m.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (m Message) FromHash(data map[string]string) (Message, error) {

	channelID, err := strconv.ParseUint(data[rediskey.REDIS_FIELD_MESSAGE_CHANNELID], 10, 64)
	if err != nil {
		return Message{}, err
	}
	userID, err := strconv.ParseUint(data[rediskey.REDIS_FIELD_MESSAGE_USERID], 10, 64)
	if err != nil {
		return Message{}, err
	}
	var replyToID *uint
	if v, ok := data[rediskey.REDIS_FIELD_MESSAGE_REPLYTOID]; ok && v != "" {
		id64, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			tmp := uint(id64)
			replyToID = &tmp
		}
	}

	createdAt, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_MESSAGE_CREATEDAT])
	if err != nil {
		return Message{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_MESSAGE_UPDATEDAT])
	if err != nil {
		return Message{}, err
	}
	return Message{
		ChannelID: uint(channelID),
		UserID:    uint(userID),

		MsgType:   data[rediskey.REDIS_FIELD_MESSAGE_MSGTYPE],
		Status:    data[rediskey.REDIS_FIELD_MESSAGE_STATUS],
		ReplyToID: replyToID,

		Content:   data[rediskey.REDIS_FIELD_MESSAGE_CONTENT],
		AttachURL: data[rediskey.REDIS_FIELD_MESSAGE_ATTACHURL],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (m Message) ToDomain(id uint) domain.Message {
	return domain.Message{
		ID:        id,
		ChannelID: m.ChannelID,
		UserID:    m.UserID,

		MsgType: m.MsgType,
		Status:  m.Status,

		ReplyToID: m.ReplyToID,

		Content:   m.Content,
		AttachURL: m.AttachURL,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m Message) FromDomain(message domain.Message) Message {
	return Message{
		ChannelID: message.ChannelID,
		UserID:    message.UserID,

		MsgType: message.MsgType,
		Status:  message.Status,

		ReplyToID: message.ReplyToID,

		Content:   message.Content,
		AttachURL: message.AttachURL,

		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}
