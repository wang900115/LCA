package gorminterface

import "github.com/wang900115/LCA/internal/domain/entities"

type UserImplement interface {
	// Create User with username and channelUUID
	CreateUser(entities.User) (entities.User, error)
	// Delete User with UserUUID
	DeleteUser(string) (entities.User, error)
}
