package repository

import (
	"LCA/internal/adapter/model"
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"
	"time"

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

func (r *MessageRepository) CreateMessage(channelUUID, userUUID, content string) (entities.Message, error) {
	messageUUID := uuid.New().String()
	createdAt := time.Now()

	message := model.Message{
		UUID:        messageUUID,
		ChannelUUID: channelUUID,
		UserUUID:    userUUID,
		Content:     content,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	if err := r.gorm.Create(&message).Error; err != nil {
		return entities.Message{}, err
	}
	return message.ToDomain(), nil
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
