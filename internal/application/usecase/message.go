package usecase

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	gorminterface "github.com/wang900115/LCA/internal/domain/interface/gorm"
)

type MessageUsecase struct {
	messageRepo gorminterface.MessageImplement
}

func NewMessageUsecase(messageRepo gorminterface.MessageImplement) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
	}
}

func (m *MessageUsecase) CreateMessage(channel, user, content string) (string, error) {
	messageDomain := entities.Message{
		Channel: channel,
		User:    user,
		Content: content,
	}

	message, err := m.messageRepo.CreateMessage(messageDomain)
	if err != nil {
		return "", err
	}
	return message.User, nil
}

func (m *MessageUsecase) DeleteMessage(messageUUID string) error {
	err := m.messageRepo.DeleteMessage(messageUUID)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageUsecase) QueryMessages(channel string) ([]string, error) {
	messages, err := m.messageRepo.QueryMessages(channel)
	if err != nil {
		return nil, err
	}
	var Content []string
	for _, message := range messages {
		Content = append(Content, message.Content)
	}
	return Content, nil
}
