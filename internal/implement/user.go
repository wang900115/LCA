package implement

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	gormmodel "github.com/wang900115/LCA/internal/adapter/gorm/model"
	rediskey "github.com/wang900115/LCA/internal/adapter/redis/key"
	"github.com/wang900115/LCA/internal/domain/entities"
	"golang.org/x/crypto/argon2"

	"gorm.io/gorm"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

type UserImplement interface {
	Create(context.Context, entities.User) error
	Read(context.Context, uint) (*entities.User, error)
	Update(context.Context, uint, string, any) error
	Delete(context.Context, uint) error
	VerifyLogin(context.Context, string, string) (*uint, error)
	VerifyLogout(context.Context, string) (*uint, error)
	CreateLogin(context.Context, uint, entities.UserLogin) error
	UpdateLogin(context.Context, uint, int64) (*entities.UserLogin, error)
	UpdateChannel(context.Context, uint, uint, int64) (*entities.UserJoin, error)
}

type UserRepository struct {
	gorm  *gorm.DB
	redis *redis.Client
}

func NewUserRepository(gorm *gorm.DB, redis *redis.Client) UserImplement {
	return &UserRepository{
		gorm:  gorm,
		redis: redis,
	}
}

func (r *UserRepository) Create(ctx context.Context, user entities.User) error {
	password, err := hashPasswordArgon2id(*user.Password)
	if err != nil {
		return err
	}
	userModel := gormmodel.User{
		Username: user.Username,
		Password: password,
		NickName: user.NickName,
		FullName: user.FullName,
		LastName: user.LastName,
		Email:    user.Email,
		Phone:    user.Phone,
		Birth:    user.Birth,
	}
	if err := r.gorm.WithContext(ctx).Create(&userModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Read(ctx context.Context, id uint) (*entities.User, error) {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Model(&user).First(id).Error; err != nil {
		return nil, err
	}
	return user.ToDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, id uint, field string, value any) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(user).Update(field, value).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	if err := r.gorm.WithContext(ctx).Unscoped().Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) VerifyLogin(ctx context.Context, username string, password string) (*uint, error) {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	pass, err := verifyPasswordArgon2id(password, user.Password)
	if err != nil {
		return nil, err
	}
	if !pass {
		err := fmt.Errorf("password invalid")
		return nil, err
	}
	return &user.ToDomain().ID, nil
}

func (r *UserRepository) VerifyLogout(ctx context.Context, ipAddress string) (*uint, error) {
	var userLogin gormmodel.UserLogin
	if err := r.gorm.WithContext(ctx).Where("ip_address = ?", ipAddress).First(&userLogin).Error; err != nil {
		return nil, err
	}
	return &userLogin.UserID, nil
}

func (r *UserRepository) CreateLogin(ctx context.Context, id uint, userLogin entities.UserLogin) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	userLoginModel := gormmodel.UserLogin{
		UserID:     id,
		LastLogin:  userLogin.LastLogin,
		IPAddress:  *userLogin.IPAddress,
		DeviceInfo: *userLogin.DeviceInfo,
		User:       user,
	}

	if err := r.gorm.WithContext(ctx).Create(&userLoginModel).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_LOGIN_TABLE + strconv.FormatUint(uint64(userLoginModel.ID), 10)
	fields := map[string]interface{}{
		rediskey.REDIS_USER_LOGIN_FIELD_IPADDRESS:  userLogin.IPAddress,
		rediskey.REDIS_USER_LOGIN_FIELD_DEVICEINFO: userLogin.DeviceInfo,
		rediskey.REDIS_USER_LOGIN_FIELD_LASTLOGIN:  userLogin.LastLogin,
	}

	go func(ctx context.Context, key string, fields map[string]interface{}) {
		if err := r.redis.HSet(ctx, key, fields).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 7*24*time.Hour)
	}(ctx, key, fields)
	return nil
}

func (r *UserRepository) UpdateLogin(ctx context.Context, login_id uint, loginTime int64) (*entities.UserLogin, error) {
	var userLogin gormmodel.UserLogin
	if err := r.gorm.WithContext(ctx).Where("id = ?", login_id).First(&userLogin).Update("last_login", loginTime).Error; err != nil {
		return nil, err
	}

	key := rediskey.REDIS_USER_LOGIN_TABLE + strconv.FormatUint(uint64(login_id), 10)
	go func(ctx context.Context, key string, lastLogin int64) {
		if err := r.redis.HSet(ctx, key, rediskey.REDIS_USER_LOGIN_FIELD_LASTLOGIN, lastLogin).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 7*24*time.Hour)
	}(ctx, key, loginTime)
	return userLogin.ToDomain(), nil
}

func (r *UserRepository) CreateChannel(ctx context.Context, id uint, channel_id uint, userJoin entities.UserJoin) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	var channel gormmodel.Channel
	if err := r.gorm.WithContext(ctx).Where("id = ?", channel_id).First(&channel).Error; err != nil {
		return err
	}

	userJoinModel := gormmodel.MiddleChannelUser{
		UserID:    id,
		ChannelID: id,
		Role:      userJoin.Role,
		LastJoin:  userJoin.LastJoin,
	}

	if err := r.gorm.WithContext(ctx).Create(&userJoinModel).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_CHANNEL_TABLE + strconv.FormatUint(uint64(userJoinModel.ID), 10)
	fields := map[string]interface{}{
		rediskey.REDIS_USER_CHANNEL_FIELD_USERID:    id,
		rediskey.REDIS_USER_CHANNEL_FIELD_CHANNELID: channel_id,
		rediskey.REDIS_USER_CHANNEL_FIELD_LASTJOIN:  userJoin.LastJoin,
		rediskey.REDIS_USER_CHANNEL_FIELD_ROLE:      userJoin.Role,
	}

	go func(ctx context.Context, key string, fields map[string]interface{}) {
		if err := r.redis.HSet(ctx, key, fields).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 24*time.Hour)
	}(ctx, key, fields)
	return nil
}

func (r *UserRepository) UpdateChannel(ctx context.Context, id uint, channel_id uint, joinTime int64) (*entities.UserJoin, error) {
	var userChannel gormmodel.MiddleChannelUser
	if err := r.gorm.WithContext(ctx).Where("user_id = ? && channel_id =?", id, channel_id).First(&userChannel).Update("last_join", joinTime).Error; err != nil {
		return nil, err
	}

	key := rediskey.REDIS_USER_CHANNEL_TABLE + strconv.FormatUint(uint64(userChannel.ID), 10)
	go func(ctx context.Context, key string, lastLogin int64) {
		if err := r.redis.HSet(ctx, key, rediskey.REDIS_USER_CHANNEL_FIELD_LASTJOIN, lastLogin).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 7*24*time.Hour)
	}(ctx, key, joinTime)
	return userChannel.ToDomain(), nil
}

func hashPasswordArgon2id(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.Key([]byte(password), salt, iterations, memory, parallelism, keyLength)
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hash)
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64salt, b64hash)
	return encoded, nil
}

func verifyPasswordArgon2id(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, errors.New("invalid hash format")
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
