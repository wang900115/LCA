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

func (u *UserUsecase) CreateUser(channelname, username string) (string, error) {
	userDomain := entities.User{
		Username: username,
		Channel:  channelname,
	}
	user, err := u.userRepo.CreateUser(userDomain)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func (u *UserUsecase) DeleteUser(username string) (string, error) {
	user, err := u.userRepo.DeleteUser(username)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}
