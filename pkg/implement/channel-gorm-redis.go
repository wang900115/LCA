package implement

import (
	"context"
	"strconv"
	"time"

	"github.com/wang900115/LCA/pkg/bootstrap"
	"github.com/wang900115/LCA/pkg/domain"
	gormmodel "github.com/wang900115/LCA/pkg/gorm/model"
	rediskey "github.com/wang900115/LCA/pkg/redis/key"
	redismodel "github.com/wang900115/LCA/pkg/redis/model"
	"go.uber.org/zap"
)

type ChannelQueryService interface {
	// 查詢所有的頻道
	QueryChannel(c context.Context) ([]domain.Channel, error)
	// 查詢用戶的頻道
	QueryCertainChannel(c context.Context, user_id uint) ([]domain.Channel, error)
	// 查詢該頻道的用戶
	QueryUser(c context.Context, channel_id uint) ([]domain.User, error)
}

type ChannelCommandService interface {
	// 創建頻道
	CreateChannel(c context.Context, toCreate domain.Channel) (domain.Channel, error)
	// 刪除特定頻道
	DeleteChannel(c context.Context, channel_id uint) error
	// 更新頻道
	UpdateChannel(c context.Context, toUpdate domain.Channel) (domain.Channel, error)
}

type ChannelReadRepository struct {
	gorm   *bootstrap.DBGroup
	redis  *bootstrap.RedisGroup
	logger *zap.Logger
}

type ChannelWriteRepository struct {
	gorm   *bootstrap.DBGroup
	redis  *bootstrap.RedisGroup
	logger *zap.Logger
}

func NewChannelReadRepository(gorm *bootstrap.DBGroup, redis *bootstrap.RedisGroup, logger *zap.Logger) ChannelQueryService {
	return &ChannelReadRepository{
		gorm:   gorm,
		redis:  redis,
		logger: logger,
	}
}

func NewChannelWriteRepository(gorm *bootstrap.DBGroup, redis *bootstrap.RedisGroup, logger *zap.Logger) ChannelCommandService {
	return &ChannelWriteRepository{
		gorm:   gorm,
		redis:  redis,
		logger: logger,
	}
}

func (cr *ChannelReadRepository) QueryChannel(c context.Context) ([]domain.Channel, error) {
	// 先去尋找 Read Redis 是否有資料
	var cursor uint64
	var keys []string
	pattern := rediskey.REDIS_TABLE_CHANNEL
	redisReader := cr.redis.PickRedisLeastConnRead()
	for {
		var err error
		var scannedKeys []string
		scannedKeys, cursor, err = redisReader.Scan(c, cursor, pattern+"*", 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, scannedKeys...)
		if cursor == 0 {
			break
		}
	}

	// 判斷是否有掃到
	if len(keys) > 0 {
		var channels []domain.Channel
		for _, key := range keys {
			channelKey := pattern + key
			data, err := redisReader.HGetAll(c, channelKey).Result()
			if err != nil {
				return nil, err
			}
			redisModel, err := redismodel.Channel{}.FromHash(data)
			if err != nil {
				return nil, err
			}
			id, err := strconv.ParseUint(key, 10, 64)
			if err != nil {
				return nil, err
			}
			channels = append(channels, redisModel.ToDomain(uint(id)))
		}
		return channels, nil
	}

	// 從資料庫查找
	var channelsModel []gormmodel.Channel
	var channels []domain.Channel
	gormReader := cr.gorm.PickDBLeastConnRead()
	if err := gormReader.WithContext(c).Find(&channelsModel).Error; err != nil {
		return nil, err
	}

	for _, channelModel := range channelsModel {
		channels = append(channels, channelModel.ToDomain())
	}

	// 背景更新 redis
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, domains []domain.Channel) {
		for _, channel := range domains {
			select {
			case <-ctx.Done():
				return
			default:
				model := redismodel.Channel{}.FromDomain(channel)
				tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(channel.ID))
				if err := cr.redis.Write.Redis.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
					cr.logger.Error("Redis Write Channel Table Err", zap.Error(err))
				}
			}
		}
	}(ctx, channels)

	return channels, nil
}

func (cr *ChannelReadRepository) QueryCertainChannel(c context.Context, user_id uint) ([]domain.Channel, error) {
	// 從 redis 取得所有channel_id 再取 table
	setKey := rediskey.REDIS_SET_USER_CHANNELS + strconv.FormatUint(uint64(user_id), 10)
	redisReader := cr.redis.PickRedisLeastConnRead()
	channelIDs, err := redisReader.SMembers(c, setKey).Result()
	if err == nil && len(channelIDs) > 0 {
		var channels []domain.Channel
		for _, channelID := range channelIDs {
			tableKey := rediskey.REDIS_TABLE_CHANNEL + channelID
			data, err := redisReader.HGetAll(c, tableKey).Result()
			if err != nil || len(data) == 0 {
				return nil, err
			}
			redisModel, err := redismodel.Channel{}.FromHash(data)
			if err != nil {
				return nil, err
			}
			id, err := strconv.ParseUint(channelID, 10, 64)
			if err != nil {
				return nil, err
			}
			channels = append(channels, redisModel.ToDomain(uint(id)))
		}
		return channels, nil
	}

	// 從 database 取
	var channelModels []gormmodel.Channel
	var channels []domain.Channel
	gormReader := cr.gorm.PickDBLeastConnRead()
	err = gormReader.WithContext(c).Joins("JOIN user_channels ON user_channels.channel_id = channels.id").Where("user_channels.user_id = ?", user_id).Find(&channelModels).Error
	if err != nil {
		return nil, err
	}
	for _, channelModel := range channelModels {
		channels = append(channels, channelModel.ToDomain())
	}

	// 背景存入 redis set
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, domains []domain.Channel) {
		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cr.redis.Write.Redis.SAdd(ctx, setKey, domain.ID).Err(); err != nil {
					cr.logger.Error("Redis Write User-Channels Set Err ", zap.Error(err))
				}
				tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(domain.ID))
				if err := cr.redis.Write.Redis.HSet(ctx, tableKey, redismodel.Channel{}.FromDomain(domain).ToHash()).Err(); err != nil {
					cr.logger.Error("Redis Write Channel Table Err ", zap.Error(err))
				}
			}
		}
	}(ctx, channels)

	return channels, nil
}

func (cr *ChannelReadRepository) QueryUser(c context.Context, channel_id uint) ([]domain.User, error) {
	setKey := rediskey.REDIS_SET_CHANNEL_USER + strconv.FormatUint(uint64(channel_id), 10)
	redisReader := cr.redis.PickRedisLeastConnRead()
	userIDs, err := redisReader.SMembers(c, setKey).Result()
	if err == nil && len(userIDs) > 0 {
		var users []domain.User
		for _, userID := range userIDs {
			tableKey := rediskey.REDIS_TABLE_USER + userID
			data, err := redisReader.HGetAll(c, tableKey).Result()
			if err != nil || len(data) == 0 {
				return nil, err
			}
			redisModel, err := redismodel.User{}.FromHash(data)
			if err != nil {
				return nil, err
			}
			id, err := strconv.ParseUint(userID, 10, 64)
			if err != nil {
				return nil, err
			}
			users = append(users, redisModel.ToDomain(uint(id)))
		}
		return users, nil
	}

	// 從 database 取
	var userModels []gormmodel.User
	var users []domain.User
	gormReader := cr.gorm.PickDBLeastConnRead()
	err = gormReader.WithContext(c).Joins("JOIN channel_users ON channel_users.user_id = users.id").Where("channel_users.channel_id = ?", channel_id).Find(&userModels).Error
	if err != nil {
		return nil, err
	}
	for _, userModel := range userModels {
		users = append(users, userModel.ToDomain())
	}

	// 背景存入 redis set
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, domains []domain.User) {
		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return
			default:
				if err := cr.redis.Write.Redis.SAdd(ctx, setKey, domain.ID).Err(); err != nil {
					cr.logger.Error("Redis Write Channel-Users Set Err ", zap.Error(err))
				}
				tableKey := rediskey.REDIS_TABLE_USER + strconv.Itoa(int(domain.ID))
				if err := cr.redis.Write.Redis.HSet(ctx, tableKey, redismodel.User{}.FromDomain(domain).ToHash()).Err(); err != nil {
					cr.logger.Error("Redis Write User Table Err ", zap.Error(err))
				}
			}
		}
	}(ctx, users)

	return users, nil
}

func (cw *ChannelWriteRepository) CreateChannel(c context.Context, toCreate domain.Channel) (domain.Channel, error) {
	// 先在 database 創建
	createdModel := gormmodel.Channel{}.FromDomain(toCreate)
	if err := cw.gorm.Write.DB.WithContext(c).Create(&createdModel).Error; err != nil {
		return domain.Channel{}, nil
	}

	// 背景 redis 創建
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, channel domain.Channel) {
		select {
		case <-ctx.Done():
			return
		default:
			model := redismodel.Channel{}.FromDomain(channel)
			tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(channel.ID))
			if err := cw.redis.Write.Redis.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
				cw.logger.Error("Redis Write Creat Channel Table Err ", zap.Error(err))
			}
		}
	}(ctx, toCreate)

	return toCreate, nil
}

func (cw *ChannelWriteRepository) DeleteChannel(c context.Context, channel_id uint) error {
	// 先在 database 刪除
	if err := cw.gorm.Write.DB.WithContext(c).Delete(&gormmodel.Channel{}, channel_id).Error; err != nil {
		return err
	}

	// 背景 redis 刪除
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, channel_id uint) {
		select {
		case <-ctx.Done():
			return
		default:
			tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(channel_id))
			if err := cw.redis.Write.Redis.Del(ctx, tableKey).Err(); err != nil {
				cw.logger.Error("Redis Write Delete Channel Table Err ", zap.Error(err))
			}
		}
	}(ctx, channel_id)

	return nil
}

func (cw *ChannelWriteRepository) UpdateChannel(c context.Context, toUpdate domain.Channel) (domain.Channel, error) {
	// 先在 database 更新
	updatedModel := gormmodel.Channel{}.FromDomain(toUpdate)
	if err := cw.gorm.Write.DB.WithContext(c).Updates(updatedModel).Error; err != nil {
		return domain.Channel{}, err
	}

	// 背景 redis 更新
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, channel domain.Channel) {
		select {
		case <-ctx.Done():
			return
		default:
			model := redismodel.Channel{}.FromDomain(channel)
			tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(channel.ID))
			if err := cw.redis.Write.Redis.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
				cw.logger.Error("Redis Write Update Channel Table Err ", zap.Error(err))
			}
		}
	}(ctx, toUpdate)

	return toUpdate, nil
}
