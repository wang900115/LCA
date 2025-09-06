package implement

import (
	"context"

	"github.com/redis/go-redis/v9"
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"

	"gorm.io/gorm"
)

type MessageImplement interface {
	Create(context.Context, entities.Message) error
	Read(context.Context, uint) (*entities.Message, error)
	Update(context.Context, uint, string, any) error
	Delete(context.Context, uint) error
}

type MessageRepository struct {
	gorm  *gorm.DB
	redis *redis.Client
}

func NewMessageRepository(gorm *gorm.DB, redis *redis.Client) MessageImplement {
	return &MessageRepository{
		gorm:  gorm,
		redis: redis,
	}
}

func (r *MessageRepository) Create(ctx context.Context, message entities.Message) error {
	messageModel := gormmodel.Message{
		Sender:  message.Sender,
		Content: message.Content,
	}
	if err := r.gorm.WithContext(ctx).Create(&messageModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *MessageRepository) Read(ctx context.Context, id uint) (*entities.Message, error) {
	var message gormmodel.Message
	if err := r.gorm.Where("id = ?", id).First(&message).Error; err != nil {
		return nil, err
	}
	return message.ToDomain(), nil
}

func (r *MessageRepository) Update(ctx context.Context, id uint, field string, value any) error {
	if err := r.gorm.WithContext(ctx).Table("message").Where("id = ?", id).Update(field, value).Error; err != nil {
		return err
	}
	return nil
}

func (r *MessageRepository) Delete(ctx context.Context, id uint) error {
	var message gormmodel.Message
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&message).Error; err != nil {
		return err
	}
	if err := r.gorm.WithContext(ctx).Unscoped().Delete(&message).Error; err != nil {
		return err
	}
	return nil
}
