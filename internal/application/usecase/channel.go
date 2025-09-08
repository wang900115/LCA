package usecase

import (
	"context"

	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/implement"
)

type ChannelUsecase struct {
	channelRepo implement.ChannelImplement
}

func NewChannelUsecase(channelRepo implement.ChannelImplement) *ChannelUsecase {
	return &ChannelUsecase{
		channelRepo: channelRepo,
	}
}

func (c *ChannelUsecase) CreateChannel(ctx context.Context, id uint, req validator.ChannelCreateRequest) error {
	channel := entities.Channel{
		Name:        req.Name,
		FounderID:   id,
		ChannelType: req.ChannelType,
	}
	return c.channelRepo.Create(ctx, id, channel)
}

func (c *ChannelUsecase) ReadChannel(ctx context.Context, id uint) (*entities.Channel, error) {
	return c.channelRepo.Read(ctx, id)
}

func (c *ChannelUsecase) UpdateChannel(ctx context.Context, id uint, field string, value any) error {
	return c.channelRepo.Update(ctx, id, field, value)
}

func (c *ChannelUsecase) DeleteChannel(ctx context.Context, id uint) error {
	return c.channelRepo.Delete(ctx, id)
}

func (c *ChannelUsecase) QueryChannelUsers(ctx context.Context, id uint) ([]*entities.User, error) {
	return c.channelRepo.ReadUsers(ctx, id)
}

func (c *ChannelUsecase) QueryChannelMessages(ctx context.Context, id uint) ([]*entities.Message, error) {
	return c.channelRepo.ReadMessages(ctx, id)
}

func (c *ChannelUsecase) UserJoin(ctx context.Context, id uint, user_id uint) error {
	return c.channelRepo.AddUser(ctx, id, user_id)
}

func (c *ChannelUsecase) UserLeave(ctx context.Context, id uint, user_id uint) error {
	return c.channelRepo.RemoveUser(ctx, id, user_id)
}

// func (c *ChannelUsecase) CommentMessage(ctx context.Context, id uint, channel_id uint, req validator.UserCommentRequest) error {
// 	message := entities.Message{
// 		Sender:  id,
// 		Content: req.Content,
// 	}
// 	return c.channelRepo.AddMessage(ctx, id, message)
// }

// func (c *ChannelUsecase) EditeMessage(ctx context.Context, id uint, message_id uint, new string) error {
// 	return c.channelRepo.UpdateMessage(ctx, id, message_id, new)
// }

// func (c *ChannelUsecase) RegainMessage(ctx context.Context, id uint, message_id uint) error {
// 	return c.channelRepo.RemoveMessage(ctx, id, message_id)
// }
