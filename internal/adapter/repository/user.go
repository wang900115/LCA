package repository

import (
	"github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/domain/irepository"

	"gorm.io/gorm"
)

type UserRepository struct {
	gorm *gorm.DB
}

func NewUserRepository(gorm *gorm.DB) irepository.IUserRepository {
	return &UserRepository{
		gorm: gorm,
	}
}

func (r *UserRepository) CreateUser(user entities.User) (entities.User, error) {
	userModel := model.User{
		Username:    user.Username,
		ChannelName: user.Channel,
	}

	if err := r.gorm.Create(&userModel).Error; err != nil {
		return entities.User{}, err
	}

	return userModel.ToDomain(), nil
}

func (r *UserRepository) DeleteUser(userName string) (entities.User, error) {
	var user model.User
	if err := r.gorm.Where("username = ?", userName).First(&user).Error; err != nil {
		return entities.User{}, err
	}

	if err := r.gorm.Unscoped().Delete(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user.ToDomain(), nil
}
