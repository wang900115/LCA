package gormimplement

import (
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"
	gorminterface "github.com/wang900115/LCA/internal/domain/interface/gorm"

	"gorm.io/gorm"
)

type UserRepository struct {
	gorm *gorm.DB
}

func NewUserRepository(gorm *gorm.DB) gorminterface.UserImplement {
	return &UserRepository{
		gorm: gorm,
	}
}

func (r *UserRepository) CreateUser(user entities.User) (entities.User, error) {
	userModel := gormmodel.User{
		Username:    user.Username,
		ChannelName: user.Channel,
	}

	if err := r.gorm.Create(&userModel).Error; err != nil {
		return entities.User{}, err
	}

	return userModel.ToDomain(), nil
}

func (r *UserRepository) DeleteUser(userName string) (entities.User, error) {
	var user gormmodel.User
	if err := r.gorm.Where("username = ?", userName).First(&user).Error; err != nil {
		return entities.User{}, err
	}

	if err := r.gorm.Unscoped().Delete(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user.ToDomain(), nil
}
