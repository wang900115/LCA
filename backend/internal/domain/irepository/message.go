package irepository

import "LCA/internal/domain/entities"

type IMessageRepository interface {
	// Create Message
	CreateMessage(entities.Message) (entities.Message, error)
	// Delete Message with messageUUID
	DeleteMessage(string) error
	// Get Message with ChannelUUID
	QueryMessages(string) ([]entities.Message, error)
}
