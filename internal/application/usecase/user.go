package usecase

import (
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/domain/irepository"
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
		Username:    username,
		ChannelUUID: channelUUID,
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
