package irepository

import "LCA/internal/domain/entities"

type IMessageRepository interface {
	// Create Message with channelUUID and userUUID and content
	CreateMessage(string, string, string) (entities.Message, error)
	// Delete Message with messageUUID
	DeleteMessage(string) error
	// Get Message with ChannelUUID
	QueryMessages(string) ([]entities.Message, error)
}
