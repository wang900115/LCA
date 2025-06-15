package repository

import (
	"github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/domain/irepository"

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

func (r *ChannelRepository) CreateChannel(name string) (entities.Channel, error) {
	channel := model.Channel{
		Name: name,
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

func (r *ChannelRepository) QueryUsers(channelName string) ([]entities.User, error) {
	var channel model.Channel
	if err := r.gorm.Preload("Users").Where("name = ?", channelName).First(&channel).Error; err != nil {
		return nil, err
	}

	var users []entities.User
	for _, user := range channel.Users {
		users = append(users, user.ToDomain())
	}
	return users, nil
}
