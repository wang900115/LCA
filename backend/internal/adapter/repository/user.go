package repository

import (
	"LCA/internal/adapter/model"
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"
	"time"

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

func (r *UserRepository) CreateUser(username, channelUUID string) (entities.User, error) {
	userUUID := uuid.New().String()
	createdAt := time.Now()
	user := model.User{
		UUID:        userUUID,
		Username:    username,
		ChannelUUID: channelUUID,
		Status:      "online",
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	if err := r.gorm.Create(&user).Error; err != nil {
		return entities.User{}, err
	}
	return user.ToDomain(), nil
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
