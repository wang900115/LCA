package implement

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wang900115/LCA/pkg/bootstrap"
	"github.com/wang900115/LCA/pkg/common"
	"github.com/wang900115/LCA/pkg/constant"
	"github.com/wang900115/LCA/pkg/domain"
	gormmodel "github.com/wang900115/LCA/pkg/gorm/model"
	rediskey "github.com/wang900115/LCA/pkg/redis/key"
	redismodel "github.com/wang900115/LCA/pkg/redis/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"
)

type UserQueryService interface {
	// 列出所有使用者
	QueryUser(c context.Context) ([]domain.User, error)

	// 檢查帳號密碼
	CheckPassword(c context.Context, username string, password string) (domain.User, error)

	// // 登入 (會包含token)
	// Login(username string, password string) error
	// // 登出 (會包含token)
	// Logout(user_id uint) error
}

type UserCommandService interface {
	// 創建使用者
	CreateUser(c context.Context, user domain.User) (domain.User, error)
	// 刪除使用者
	DeleteUser(c context.Context, user_id uint) error
	// 更新使用者
	UpdateUser(c context.Context, toUpdate domain.User) (domain.User, error)
}

type UserReadRepository struct {
	gorm   *bootstrap.DBGroup
	redis  *bootstrap.RedisGroup
	logger *zap.Logger
}

type UserWriteRepository struct {
	gorm   *bootstrap.DBGroup
	redis  *bootstrap.RedisGroup
	logger *zap.Logger
}

func NewUserReadRepository(gorm *bootstrap.DBGroup, redis *bootstrap.RedisGroup, logger *zap.Logger) UserQueryService {
	return &UserReadRepository{
		gorm:   gorm,
		redis:  redis,
		logger: logger,
	}
}

func NewUserWriteRepository(gorm *bootstrap.DBGroup, redis *bootstrap.RedisGroup, logger *zap.Logger) UserCommandService {
	return &UserWriteRepository{
		gorm:   gorm,
		redis:  redis,
		logger: logger,
	}
}

func (ur *UserReadRepository) QueryUser(c context.Context) ([]domain.User, error) {
	// 先去尋找 Read Redis 是否有資料
	var cursor uint64
	var keys []string
	pattern := rediskey.REDIS_TABLE_USER
	redisReader := ur.redis.PickRedisLeastConnRead()
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
		var users []domain.User
		for _, key := range keys {
			tableKey := pattern + key
			data, err := redisReader.HGetAll(c, tableKey).Result()
			if err != nil {
				return nil, err
			}
			redisModel, err := redismodel.User{}.FromHash(data)
			if err != nil {
				return nil, err
			}
			id, err := strconv.ParseUint(key, 10, 64)
			if err != nil {
				return nil, err
			}
			users = append(users, redisModel.ToDomain(uint(id)))
		}
		return users, nil
	}

	// 從資料庫查找
	var usersModel []gormmodel.User
	var users []domain.User
	gormReader := ur.gorm.PickDBLeastConnRead()
	if err := gormReader.WithContext(c).Find(&usersModel).Error; err != nil {
		return nil, err
	}

	for _, userModel := range usersModel {
		users = append(users, userModel.ToDomain())
	}

	// 背景更新 redis
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, domains []domain.User) {
		for _, user := range domains {
			select {
			case <-ctx.Done():
				return
			default:
				model := redismodel.User{}.FromDomain(user)
				tableKey := rediskey.REDIS_TABLE_CHANNEL + strconv.Itoa(int(user.ID))
				if err := ur.redis.Write.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
					ur.logger.Error("Redis Write User Table Err ", zap.Error(err))
				}
			}
		}
	}(ctx, users)

	return users, nil
}

func (uw *UserWriteRepository) CreateUser(c context.Context, toCreate domain.User) (domain.User, error) {
	// 先在 database 創建
	createdModel := gormmodel.User{}.FromDomain(toCreate)
	hashedPassword, err := hashPasswordArgon2id(createdModel.Password)
	if err != nil {
		return domain.User{}, err
	}

	createdModel.Password = hashedPassword
	if err := uw.gorm.Write.WithContext(c).Create(&createdModel).Error; err != nil {
		return domain.User{}, nil
	}

	// 背景 redis 創建
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, user domain.User) {
		select {
		case <-ctx.Done():
			return
		default:
			model := redismodel.User{}.FromDomain(user)
			tableKey := rediskey.REDIS_TABLE_USER + strconv.Itoa(int(user.ID))
			if err := uw.redis.Write.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
				uw.logger.Error("Redis Write Creat User Table Err: ", zap.Error(err))
			}
		}
	}(ctx, toCreate)

	return toCreate, nil
}

func (uw *UserWriteRepository) DeleteUser(c context.Context, user_id uint) error {
	// 先在 database 刪除
	if err := uw.gorm.Write.WithContext(c).Delete(&gormmodel.User{}, user_id).Error; err != nil {
		return err
	}

	// 背景 redis 刪除
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, user_id uint) {
		select {
		case <-ctx.Done():
			return
		default:
			tableKey := rediskey.REDIS_TABLE_USER + strconv.Itoa(int(user_id))
			if err := uw.redis.Write.Del(ctx, tableKey).Err(); err != nil {
				uw.logger.Error("Redis Write Delete User Table Err ", zap.Error(err))
			}
		}
	}(ctx, user_id)

	return nil
}

func (uw *UserWriteRepository) UpdateUser(c context.Context, toUpdate domain.User) (domain.User, error) {
	// 先在 database 更新
	updatedModel := gormmodel.User{}.FromDomain(toUpdate)
	if err := uw.gorm.Write.WithContext(c).Updates(updatedModel).Error; err != nil {
		return domain.User{}, err
	}

	// 背景 redis 更新
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, user domain.User) {
		select {
		case <-ctx.Done():
			return
		default:
			model := redismodel.User{}.FromDomain(user)
			tableKey := rediskey.REDIS_TABLE_USER + strconv.Itoa(int(user.ID))
			if err := uw.redis.Write.HSet(ctx, tableKey, model.ToHash()).Err(); err != nil {
				uw.logger.Error("Redis Write Update User Table Err ", zap.Error(err))
			}
		}
	}(ctx, toUpdate)

	return toUpdate, nil
}

func (uw *UserWriteRepository) UpdateRole(c context.Context, user_id uint, toUpdate string) (string, error) {
	// 先在 database 更新
	if err := uw.gorm.Write.WithContext(c).Model(&gormmodel.User{}).Where("id = ?", user_id).Update("role", toUpdate).Error; err != nil {
		return "", err
	}

	// 背景 redis 更新
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	go func(ctx context.Context, user_id uint, toUpdate string) {
		select {
		case <-ctx.Done():
			return
		default:
			tableKey := rediskey.REDIS_TABLE_USER + strconv.Itoa(int(user_id))
			fieldKey := rediskey.REDIS_FIELD_USER_ROLE
			// !todo 用 redis-writer's logger
			if err := uw.redis.Write.HSet(ctx, tableKey, fieldKey, toUpdate).Err(); err != nil {
				uw.logger.Error("Redis Write Update User Role Field Err ", zap.Error(err))
			}
		}
	}(ctx, user_id, toUpdate)

	return toUpdate, nil
}

func (ur *UserReadRepository) CheckPassword(c context.Context, username string, password string) (domain.User, error) {
	var user gormmodel.User
	gormReader := ur.gorm.PickDBLeastConnRead()
	if err := gormReader.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		return domain.User{}, err
	}
	pass, err := verifyPasswordArgon2id(user.Password, password)
	if err != nil {
		return domain.User{}, err
	}
	if !pass {
		return domain.User{}, common.PasswordIncorrect
	}
	return user.ToDomain(), nil
}

func hashPasswordArgon2id(password string) (string, error) {
	salt := make([]byte, constant.SALTLENGTH)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.Key([]byte(password), salt, constant.ITERATIONS, constant.MEMORY, constant.PARALLELISM, constant.KEYLENGTH)
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, constant.MEMORY, constant.ITERATIONS, constant.PARALLELISM, b64salt, b64hash)

	return encoded, nil
}

func verifyPasswordArgon2id(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, common.HashPassword
	}
	var version, m, t, p int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, err
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	result := argon2.IDKey([]byte(password), salt, uint32(t), uint32(m), uint8(p), uint32(len(hash)))
	if subtle.ConstantTimeCompare(hash, result) == 1 {
		return true, nil
	}
	return false, nil
}
