package repository

import (
	"github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/domain/irepository"

	"github.com/google/uuid"
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
	userUUID := uuid.New().String()
	userModel := model.User{
		UUID:        userUUID,
		Username:    user.Username,
		ChannelUUID: user.ChannelUUID,
	}

	if err := r.gorm.Create(&userModel).Error; err != nil {
		return entities.User{}, err
	}

	return userModel.ToDomain(), nil
}

func (r *UserRepository) DeleteUser(userUUID string) (entities.User, error) {
	var user model.User
	if err := r.gorm.Where("uuid = ?", userUUID).First(&user).Error; err != nil {
		return entities.User{}, err
	}

	if err := r.gorm.Unscoped().Delete(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user.ToDomain(), nil
}
