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

type MessageQueryService interface {
	// 查詢該頻道訊息
	QueryMessage(c context.Context, channel_id uint) ([]domain.Message, error)
	// 查詢該頻道特定用戶的訊息
	QueryCertainMessage(c context.Context, channel_id uint, user_id uint) ([]domain.Message, error)
}

type MessageCommandService interface {
	// 紀錄訊息
	CreateMessage(c context.Context, toCreate domain.Message) (domain.Message, error)
	// 刪除訊息
	DeleteMessage(c context.Context, message_id uint) error
	// 更新訊息
	UpdateMessage(c context.Context, toUpdate domain.Message) (domain.Message, error)
}

type MessageReadRepository struct {
	gorm   *bootstrap.DBGroup
	redis  *bootstrap.RedisGroup
	logger *zap.Logger
}

type MessageWriteRepository struct {
	gorm   *bootstrap.DBGroup
	redis  *bootstrap.RedisGroup
	logger *zap.Logger
}

func NewMessageReadRepository(gorm *bootstrap.DBGroup, redis *bootstrap.RedisGroup, logger *zap.Logger) MessageQueryService {
	return &MessageReadRepository{
		gorm:   gorm,
		redis:  redis,
		logger: logger,
	}
}

func NewMessageWriteRepository(gorm *bootstrap.DBGroup, redis *bootstrap.RedisGroup, logger *zap.Logger) MessageCommandService {
	return &MessageWriteRepository{
		gorm:   gorm,
		redis:  redis,
		logger: logger,
	}
}

func (mr *MessageReadRepository) QueryMessage(c context.Context, channel_id uint) ([]domain.Message, error) {
	// 先從 redis 查找
	setKey := rediskey.REDIS_SET_CHANNEL_MESSAGE + strconv.FormatUint(uint64(channel_id), 10)
	redisReader := mr.redis.PickRedisLeastConnRead()
	messageIDs, err := redisReader.SMembers(c, setKey).Result()
	if err == nil && len(messageIDs) > 0 {
		var messages []domain.Message
		for _, messageID := range messageIDs {
			tableKey := rediskey.REDIS_TABLE_MESSAGE + messageID
			data, err := redisReader.HGetAll(c, tableKey).Result()
			if err != nil || len(data) == 0 {
				return nil, err
			}
			redisModel, err := redismodel.Message{}.FromHash(data)
			if err != nil {
				return nil, err
			}
			id, err := strconv.ParseUint(messageID, 10, 64)
			if err != nil {
				return nil, err
			}
			messages = append(messages, redisModel.ToDomain(uint(id)))
		}
		return messages, nil
	}
	// 從 database
	var messageModels []gormmodel.Message
	var messages []domain.Message
	gormReader := mr.gorm.PickDBLeastConnRead()
	err = gormReader.WithContext(c).Joins("JOIN channnel_messages ON channel_messages.message_id = messages.id").Where("channel_messages.channel_id = ?", channel_id).Find(&messageModels).Error
	if err != nil {
		return nil, err
	}
	for _, messageModel := range messageModels {
		messages = append(messages, messageModel.ToDomain())
	}

	// 背景存入 redis set
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, domains []domain.Message) {
		for _, domain := range domains {
			select {
			case <-ctx.Done():
				return
			default:
				if err := mr.redis.Write.SAdd(ctx, setKey, domain.ID).Err(); err != nil {
					mr.logger.Error("Redis Write Channel-Messages Set Err ", zap.Error(err))
				}
				tableKey := rediskey.REDIS_TABLE_MESSAGE + strconv.Itoa(int(domain.ID))
				if err := mr.redis.Write.HSet(ctx, tableKey, redismodel.Message{}.FromDomain(domain).ToHash()).Err(); err != nil {
					mr.logger.Error("Redis Write Message Table Err ", zap.Error(err))
				}
			}
		}
	}(ctx, messages)

	return messages, nil

}

func (mr *MessageReadRepository) QueryCertainMessage(c context.Context, channel_id uint, user_id uint) ([]domain.Message, error) {
	// 從 redis 取得該 頻道特定用戶的 message_id
	listKey := rediskey.REDIS_LIST_CHANNEL_USER_MESSAGE + strconv.FormatUint(uint64(channel_id), 10) + strconv.FormatUint(uint64(user_id), 10)
	redisReader := mr.redis.PickRedisLeastConnRead()
	msgIDs, err := redisReader.LRange(c, listKey, 0, -1).Result()
	if err == nil && len(msgIDs) > 0 {
		var result []domain.Message
		for _, msgID := range msgIDs {
			messageKey := rediskey.REDIS_TABLE_MESSAGE + msgID
			data, err := redisReader.HGetAll(c, messageKey).Result()
			if err != nil {
				return nil, err
			}
			redisModel, err := redismodel.Message{}.FromHash(data)
			if err != nil {
				return nil, err
			}
			id, err := strconv.ParseUint(msgID, 10, 64)
			if err != nil {
				return nil, err
			}
			domainMsg := redisModel.ToDomain(uint(id))
			result = append(result, domainMsg)
		}
		return result, nil
	}
	// 從 database 取得
	var messageModels []gormmodel.Message
	var messages []domain.Message
	gormReader := mr.gorm.PickDBLeastConnRead()
	err = gormReader.WithContext(c).Where("channel_id = ? AND user_id = ?", channel_id, user_id).Find(&messageModels).Error
	if err != nil {
		return nil, err
	}

	for _, messageModel := range messageModels {
		messages = append(messages, messageModel.ToDomain())
	}

	// 背景存入 redis list
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, messages []domain.Message) {
		for _, message := range messages {
			select {
			case <-ctx.Done():
				return
			default:
				if err := mr.redis.Write.RPush(ctx, listKey, message.ID).Err(); err != nil {
					mr.logger.Error("Redis Push Channel-User-Messages List Err ", zap.Error(err))
				}
				tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(message.ID))
				if err := mr.redis.Write.HSet(ctx, tableKey, redismodel.Message{}.FromDomain(message).ToHash()).Err(); err != nil {
					mr.logger.Error("Redis Write Message Table Err ", zap.Error(err))
				}
			}
		}
	}(ctx, messages)

	return messages, nil
}

func (mw *MessageWriteRepository) CreateMessage(c context.Context, toCreate domain.Message) (domain.Message, error) {
	// 先在 database 創建
	createdModel := gormmodel.Message{}.FromDomain(toCreate)
	if err := mw.gorm.Write.WithContext(c).Create(&createdModel).Error; err != nil {
		return domain.Message{}, nil
	}

	// 背景 redis 創建
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, message domain.Message) {
		select {
		case <-ctx.Done():
			return
		default:
			model := redismodel.Message{}.FromDomain(message)
			tableKey := rediskey.REDIS_TABLE_MESSAGE + strconv.Itoa(int(message.ID))
			if err := mw.redis.Write.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
				mw.logger.Error("Redis Write Creat Message Table Err ", zap.Error(err))
			}
		}
	}(ctx, toCreate)

	return toCreate, nil
}

func (mw *MessageWriteRepository) DeleteMessage(c context.Context, message_id uint) error {
	// 先在 database 刪除
	if err := mw.gorm.Write.WithContext(c).Delete(&gormmodel.Message{}, message_id).Error; err != nil {
		return err
	}

	// 背景 redis 刪除
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, message_id uint) {
		select {
		case <-ctx.Done():
			return
		default:
			tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(message_id))
			if err := mw.redis.Write.Del(ctx, tableKey).Err(); err != nil {
				mw.logger.Error("Redis Write Delete Message Table Err ", zap.Error(err))
			}
		}
	}(ctx, message_id)

	return nil
}

func (mw *MessageWriteRepository) UpdateMessage(c context.Context, toUpdate domain.Message) (domain.Message, error) {
	// 先在 database 更新
	updatedModel := gormmodel.Message{}.FromDomain(toUpdate)
	if err := mw.gorm.Write.WithContext(c).Updates(updatedModel).Error; err != nil {
		return domain.Message{}, err
	}

	// 背景 redis 更新
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, message domain.Message) {
		select {
		case <-ctx.Done():
			return
		default:
			model := redismodel.Message{}.FromDomain(message)
			tableKey := rediskey.REDIS_TABLE_MESSAGE + strconv.Itoa(int(message.ID))
			if err := mw.redis.Write.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
				mw.logger.Error("Redis Write Update Message Table Err ", zap.Error(err))
			}
		}
	}(ctx, toUpdate)

	return toUpdate, nil
}
