package repository

import (
	"LCA/internal/adapter/gorm/model"
	"LCA/internal/domain/entities"
	"LCA/internal/domain/irepository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelRepository struct {
	gorm *gorm.DB
}

func NewChannelRepository(gorm *gorm.DB) irepository.IChannelRepository {
	return &ChannelRepository{
		gorm: gorm,
	}
}

func (r *ChannelRepository) CreateChannel() (entities.Channel, error) {
	channelUUID := uuid.New().String()
	channel := model.Channel{
		UUID: channelUUID,
	}
	if err := r.gorm.Create(&channel).Error; err != nil {
		return entities.Channel{}, err
	}
	return channel.ToDomain(), nil
}

func (r *ChannelRepository) QueryChannels() ([]entities.Channel, error) {
	var channels []model.Channel
	if err := r.gorm.Find(&channels).Error; err != nil {
		return nil, err
	}

	var result []entities.Channel
	for _, channel := range channels {
		result = append(result, channel.ToDomain())
	}
	return result, nil
}

func (r *ChannelRepository) QueryUsers(channelUUID string) ([]entities.User, error) {
	var channel model.Channel
	if err := r.gorm.Preload("Users").Where("uuid = ?", channelUUID).First(&channel).Error; err != nil {
		return nil, err
	}

	var users []entities.User
	for _, user := range channel.Users {
		users = append(users, user.ToDomain())
	}
	return users, nil
}
