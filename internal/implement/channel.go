package implement

import (
	"context"

	"github.com/redis/go-redis/v9"
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type ChannelImplement interface {
	Create(context.Context, entities.Channel) error
	Read(context.Context, uint) (*entities.Channel, error)
	Update(context.Context, uint, string, any) error
	Delete(context.Context, uint) error
	ReadUsers(context.Context, uint) (*entities.ChannelUser, error)
	ReadMessages(context.Context, uint) (*entities.ChannelMessage, error)
	UpdateUsers(context.Context, uint, entities.User) error
	UpdateMessages(context.Context, uint, entities.Message) error
}

type ChannelRepository struct {
	gorm  *gorm.DB
	redis *redis.Client
}

func NewChannelRepository(gorm *gorm.DB, redis *redis.Client) ChannelImplement {
	return &ChannelRepository{
		gorm:  gorm,
		redis: redis,
	}
}

func (r *ChannelRepository) Create(ctx context.Context, channel entities.Channel) error {
	channelModel := gormmodel.Channel{
		Name:        channel.Name,
		Founder:     channel.Founder,
		ChannelType: channel.ChannelType,
	}
	if err := r.gorm.WithContext(ctx).Create(&channelModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChannelRepository) Read(ctx context.Context, id uint) (*entities.Channel, error) {
	var channel gormmodel.Channel
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&channel).Error; err != nil {
		return nil, err
	}
	return channel.ToDomain(), nil
}

func (r *ChannelRepository) Update(ctx context.Context, id uint, field string, value any) error {
	var channel gormmodel.Channel
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&channel).Update(field, value).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChannelRepository) Delete(ctx context.Context, id uint) error {
	var channel gormmodel.Channel
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&channel).Error; err != nil {
		return err
	}
	if err := r.gorm.WithContext(ctx).Unscoped().Delete(&channel).Error; err != nil {
		return err
	}
	return nil
}

// !todo(redis, db)
func (r *ChannelRepository) ReadUsers(ctx context.Context, id uint) (*entities.ChannelUser, error) {
	var channelUser gormmodel.ChannelUser
	if err := r.gorm.WithContext(ctx).Where("channel_id = ?", id).First(&channelUser).Error; err != nil {
		return nil, err
	}
	return channelUser.ToDomain(), nil
}

// !todo(redis, db)
func (r *ChannelRepository) UpdateUsers(ctx context.Context, id uint, user entities.User) error {
	return nil
}

// !todo(redis, db)
func (r *ChannelRepository) ReadMessages(ctx context.Context, id uint) (*entities.ChannelMessage, error) {
	var channelMessage gormmodel.ChannelMessage
	if err := r.gorm.WithContext(ctx).Where("channel_id = ?", id).First(&channelMessage).Error; err != nil {
		return nil, err
	}
	return channelMessage.ToDomain(), nil
}

// !todo(redis, db)
func (r *ChannelRepository) UpdateMessages(ctx context.Context, id uint, message entities.Message) error {
	return nil
}
