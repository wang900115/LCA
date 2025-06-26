package gorminterface

import "github.com/wang900115/LCA/internal/domain/entities"

type MessageImplement interface {
	// Create Message
	CreateMessage(entities.Message) (entities.Message, error)
	// Delete Message with messageUUID
	DeleteMessage(string) error
	// Get Message with ChannelUUID
	QueryMessages(string) ([]entities.Message, error)
}
