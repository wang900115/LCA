package implement

import (
	"context"

	"github.com/redis/go-redis/v9"
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type ChannelImplement interface {
	Create(context.Context, uint, entities.Channel) error
	Read(context.Context, uint) (*entities.Channel, error)
	Update(context.Context, uint, string, any) error
	Delete(context.Context, uint) error
	ReadUsers(context.Context, uint) ([]*entities.User, error)
	ReadMessages(context.Context, uint) ([]*entities.Message, error)
	AddUser(context.Context, uint, entities.User) error
	RemoveUser(context.Context, uint, uint) error
	AddMessage(context.Context, uint, entities.Message) error
	UpdateMessage(context.Context, uint, uint, string) error
	RemoveMessage(context.Context, uint, uint) error
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

func (r *ChannelRepository) Create(ctx context.Context, id uint, channel entities.Channel) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).First(&user, id).Error; err != nil {
		return err
	}
	channelModel := gormmodel.Channel{
		Name:        channel.Name,
		FounderID:   channel.FounderID,
		Founder:     user,
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
func (r *ChannelRepository) ReadUsers(ctx context.Context, id uint) ([]*entities.User, error) {
	return nil, nil
}

// !todo(redis, db)
func (r *ChannelRepository) ReadMessages(ctx context.Context, id uint) ([]*entities.Message, error) {
	return nil, nil
}

// !todo(redis, db)
func (r *ChannelRepository) AddUser(ctx context.Context, id uint, user entities.User) error {
	return nil
}

// !todo(redis, db)
func (r *ChannelRepository) RemoveUser(ctx context.Context, id uint, user_id uint) error {
	return nil
}

// !todo(redis, db)
func (r *ChannelRepository) AddMessage(ctx context.Context, id uint, message entities.Message) error {
	return nil
}

// !todo(redis, db)
func (r *ChannelRepository) UpdateMessage(ctx context.Context, id uint, message_id uint, new string) error {
	return nil
}

// !todo(redis, db)
func (r *ChannelRepository) RemoveMessage(ctx context.Context, id uint, message_id uint) error {
	return nil
}
