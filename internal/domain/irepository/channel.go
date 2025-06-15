package irepository

import "github.com/wang900115/LCA/internal/domain/entities"

type IChannelRepository interface {
	// Create Channel
	CreateChannel(string) (entities.Channel, error)
	// Query channel
	QueryChannels() ([]entities.Channel, error)
	// Get User with ChannelUUID
	QueryUsers(string) ([]entities.User, error)
}
