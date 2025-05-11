package usecase

import "LCA/internal/domain/irepository"

type ChannelUsecase struct {
	channelRepo irepository.IChannelRepository
}

func NewChannelUsecase(channelRepo irepository.IChannelRepository) *ChannelUsecase {
	return &ChannelUsecase{
		channelRepo: channelRepo,
	}
}

func (c *ChannelUsecase) CreateChannel() (string, error) {
	channel, err := c.channelRepo.CreateChannel()
	if err != nil {
		return "", err
	}
	return channel.UUID, nil
}

func (c *ChannelUsecase) QueryChannels() ([]string, error) {
	channels, err := c.channelRepo.QueryChannels()
	if err != nil {
		return nil, err
	}
	var uuids []string
	for _, channel := range channels {
		uuids = append(uuids, channel.UUID)
	}
	return uuids, nil
}

func (c *ChannelUsecase) QueryUsers(channelUUID string) ([]string, error) {
	users, err := c.channelRepo.QueryUsers(channelUUID)
	if err != nil {
		return nil, err
	}
	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}
	return usernames, nil
}
