package usecase

import (
	"context"

	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/implement"
)

type MessageUsecase struct {
	messageRepo implement.MessageImplement
}

func NewMessageUsecase(messageRepo implement.MessageImplement) *MessageUsecase {
	return &MessageUsecase{
		messageRepo: messageRepo,
	}
}

func (m *MessageUsecase) CreateMessage(ctx context.Context, sender uint, content string) error {
	messageDomain := entities.Message{
		Sender:  sender,
		Content: content,
	}
	return m.messageRepo.Create(ctx, messageDomain)
}

func (m *MessageUsecase) ReadMessage(ctx context.Context, id uint) (*entities.Message, error) {
	return m.messageRepo.Read(ctx, id)
}

func (m *MessageUsecase) UpdateMessage(ctx context.Context, id uint, field string, value any) error {
	return m.messageRepo.Update(ctx, id, field, value)
}

func (m *MessageUsecase) DeleteMessage(ctx context.Context, id uint) error {
	return m.messageRepo.Delete(ctx, id)
}
