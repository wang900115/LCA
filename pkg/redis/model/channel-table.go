package redismodel

import (
	"time"

	"github.com/wang900115/LCA/pkg/domain"
	rediskey "github.com/wang900115/LCA/pkg/redis/key"
)

type Channel struct {
	ChannelName string
	ChannelType string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c Channel) ToHash() map[string]interface{} {
	return map[string]interface{}{
		rediskey.REDIS_FIELD_CHANNEL_NAME:      c.ChannelName,
		rediskey.REDIS_FIELD_CHANNEL_TYPE:      c.ChannelType,
		rediskey.REDIS_FIELD_CHANNEL_CREATEDAT: c.CreatedAt.UTC().Format(time.RFC3339),
		rediskey.REDIS_FIELD_CHANNEL_UPDATEDAT: c.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (c Channel) FromHash(data map[string]string) (Channel, error) {
	createdAt, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_CHANNEL_CREATEDAT])
	if err != nil {
		return Channel{}, err
	}
	updatedAt, err := time.Parse(time.RFC3339, data[rediskey.REDIS_FIELD_CHANNEL_UPDATEDAT])
	if err != nil {
		return Channel{}, err
	}
	return Channel{
		ChannelName: data[rediskey.REDIS_FIELD_CHANNEL_NAME],
		ChannelType: data[rediskey.REDIS_FIELD_CHANNEL_TYPE],
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (c Channel) ToDomain(id uint) domain.Channel {
	return domain.Channel{
		ID:          id,
		ChannelName: c.ChannelName,
		ChannelType: c.ChannelType,

		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (c Channel) FromDomain(channel domain.Channel) Channel {
	return Channel{
		ChannelName: c.ChannelName,
		ChannelType: c.ChannelType,

		CreatedAt: channel.CreatedAt,
		UpdatedAt: channel.UpdatedAt,
	}
}
