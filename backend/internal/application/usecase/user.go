package usecase

import (
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"
)

type UserUsecase struct {
	userRepo irepository.IUserRepository
}

func NewUserUsecase(userRepo irepository.IUserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) CreateUser(channelUUID, username string) (string, error) {
	userDomain := entities.User{
		UUID:        " ",
		Username:    username,
		ChannelUUID: channelUUID,
		Status:      "online",
	}
	user, err := u.userRepo.CreateUser(userDomain)
	if err != nil {
		return "", err
	}
	return user.UUID, nil
}

func (u *UserUsecase) DeleteUser(userUUID string) (string, error) {
	user, err := u.userRepo.DeleteUser(userUUID)
	if err != nil {
		return "", err
	}
	return user.UUID, nil
}
