package gormimplement

import (
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"
	gorminterface "github.com/wang900115/LCA/internal/domain/interface/gorm"

	"gorm.io/gorm"
)

type MessageRepository struct {
	gorm *gorm.DB
}

func NewMessageRepository(gorm *gorm.DB) gorminterface.MessageImplement {
	return &MessageRepository{
		gorm: gorm,
	}
}

func (r *MessageRepository) CreateMessage(message entities.Message) (entities.Message, error) {
	messageModel := gormmodel.Message{
		ChannelName: message.Channel,
		Username:    message.User,
		Content:     message.Content,
	}
	if err := r.gorm.Create(&messageModel).Error; err != nil {
		return entities.Message{}, err
	}
	return messageModel.ToDomain(), nil
}

func (r *MessageRepository) DeleteMessage(messageID string) error {
	return r.gorm.Where("id = ?", messageID).Delete(&gormmodel.Message{}).Error
}

func (r *MessageRepository) QueryMessages(channelName string) ([]entities.Message, error) {
	var messages []gormmodel.Message
	if err := r.gorm.Where("channel_name = ?", channelName).Find(&messages).Error; err != nil {
		return nil, err
	}

	var result []entities.Message
	for _, message := range messages {
		result = append(result, message.ToDomain())
	}
	return result, nil
}
