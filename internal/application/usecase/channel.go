package usecase

import gorminterface "github.com/wang900115/LCA/internal/domain/interface/gorm"

type ChannelUsecase struct {
	channelRepo gorminterface.ChannelImplement
}

func NewChannelUsecase(channelRepo gorminterface.ChannelImplement) *ChannelUsecase {
	return &ChannelUsecase{
		channelRepo: channelRepo,
	}
}

func (c *ChannelUsecase) CreateChannel(name string) (string, error) {
	channel, err := c.channelRepo.CreateChannel(name)
	if err != nil {
		return "", err
	}
	return channel.Name, nil
}

func (c *ChannelUsecase) QueryChannels() ([]string, error) {
	channels, err := c.channelRepo.QueryChannels()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, channel := range channels {
		names = append(names, channel.Name)
	}
	return names, nil
}

func (c *ChannelUsecase) QueryUsers(channelName string) ([]string, error) {
	users, err := c.channelRepo.QueryUsers(channelName)
	if err != nil {
		return nil, err
	}
	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}
	return usernames, nil
}
