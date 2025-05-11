package repository

import (
	"LCA/internal/adapter/gorm/model"
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	gorm *gorm.DB
}

func NewMessageRepository(gorm *gorm.DB) irepository.IMessageRepository {
	return &MessageRepository{
		gorm: gorm,
	}
}

func (r *MessageRepository) CreateMessage(message entities.Message) (entities.Message, error) {
	messageUUID := uuid.New().String()

	messageModel := model.Message{
		UUID:        messageUUID,
		ChannelUUID: message.ChannelUUID,
		UserUUID:    message.UserUUID,
		Content:     message.Content,
	}
	if err := r.gorm.Create(&messageModel).Error; err != nil {
		return entities.Message{}, err
	}
	return messageModel.ToDomain(), nil
}

func (r *MessageRepository) DeleteMessage(messageUUID string) error {
	return r.gorm.Where("uuid = ?", messageUUID).Delete(&model.Message{}).Error
}

func (r *MessageRepository) QueryMessages(channelUUID string) ([]entities.Message, error) {
	var messages []model.Message
	if err := r.gorm.Where("channel_uuid = ?", channelUUID).Find(&messages).Error; err != nil {
		return nil, err
	}

	var result []entities.Message
	for _, message := range messages {
		result = append(result, message.ToDomain())
	}
	return result, nil
}
