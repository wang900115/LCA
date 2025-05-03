package repository

import (
	"LCA/internal/adapter/model"
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"
	"time"

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
	currentTime := time.Now()
	userModel := model.User{
		UUID:        user.UUID,
		Username:    user.Username,
		ChannelUUID: user.ChannelUUID,
		Status:      user.Status,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
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

	if err := r.gorm.Delete(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user.ToDomain(), nil
}
