package irepository

import "LCA/internal/domain/entities"

type IUserRepository interface {
	// Create User with username and channelUUID
	CreateUser(string, string) (entities.User, error)
	// Delete User with UserUUID
	DeleteUser(string) (entities.User, error)
}
