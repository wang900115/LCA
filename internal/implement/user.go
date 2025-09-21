package implement

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
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
	UpdateLoginTime(context.Context, uint, int64) (*entities.UserLogin, error)
	DeleteLogin(context.Context, uint, uint) error
	CreateJoin(context.Context, uint, uint, entities.UserJoin) error
	UpdateJoinTime(context.Context, uint, uint, int64) (*entities.UserJoin, error)
	DeleteJoin(context.Context, uint, uint) error
	CreateParticate(context.Context, uint, uint, entities.UserParticate) error
	UpdateParticateTime(context.Context, uint, uint, int64) (*entities.UserParticate, error)
	DeleteParticate(context.Context, uint, uint) error
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
	key := rediskey.REDIS_USER_LOGIN_TABLE + fmt.Sprintf("%d", id)
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

func (r *UserRepository) UpdateLoginTime(ctx context.Context, id uint, loginTime int64) (*entities.UserLogin, error) {
	var userLogin gormmodel.UserLogin
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&userLogin).Update("last_login", loginTime).Error; err != nil {
		return nil, err
	}

	key := rediskey.REDIS_USER_LOGIN_TABLE + fmt.Sprintf("%d", id)
	go func(ctx context.Context, key string, lastLogin int64) {
		if err := r.redis.HSet(ctx, key, rediskey.REDIS_USER_LOGIN_FIELD_LASTLOGIN, lastLogin).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 7*24*time.Hour)
	}(ctx, key, loginTime)
	return userLogin.ToDomain(), nil
}

func (r *UserRepository) DeleteLogin(ctx context.Context, id uint, login_id uint) error {
	var userLogin gormmodel.UserLogin
	if err := r.gorm.WithContext(ctx).Where("user_id = ? && id = ?", id, login_id).Delete(&userLogin).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_LOGIN_TABLE + fmt.Sprintf("%d:%d", id, login_id)
	go func(ctx context.Context, key string) {
		if err := r.redis.Del(ctx, key).Err(); err != nil {
			log.Printf("redis Del error: %v", err)
		}
	}(ctx, key)
	return nil
}

func (r *UserRepository) CreateJoin(ctx context.Context, id uint, channel_id uint, userJoin entities.UserJoin) error {
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
		ChannelID: channel_id,
		Role:      userJoin.Role,
		LastJoin:  userJoin.LastJoin,
	}

	if err := r.gorm.WithContext(ctx).Create(&userJoinModel).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_CHANNEL_TABLE + fmt.Sprintf("%d:%d", id, channel_id)
	fields := map[string]interface{}{
		rediskey.REDIS_USER_CHANNEL_FIELD_LASTJOIN: userJoin.LastJoin,
		rediskey.REDIS_USER_CHANNEL_FIELD_ROLE:     userJoin.Role,
	}

	go func(ctx context.Context, key string, fields map[string]interface{}) {
		if err := r.redis.HSet(ctx, key, fields).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 24*time.Hour)
	}(ctx, key, fields)
	return nil
}

func (r *UserRepository) UpdateJoinTime(ctx context.Context, id uint, channel_id uint, joinTime int64) (*entities.UserJoin, error) {
	var userChannel gormmodel.MiddleChannelUser
	if err := r.gorm.WithContext(ctx).Where("user_id = ? && channel_id =?", id, channel_id).First(&userChannel).Update("last_join", joinTime).Error; err != nil {
		return nil, err
	}

	key := rediskey.REDIS_USER_CHANNEL_TABLE + fmt.Sprintf("%d:%d", id, channel_id)
	go func(ctx context.Context, key string, lastLogin int64) {
		if err := r.redis.HSet(ctx, key, rediskey.REDIS_USER_CHANNEL_FIELD_LASTJOIN, lastLogin).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 7*24*time.Hour)
	}(ctx, key, joinTime)
	return userChannel.ToDomain(), nil
}

func (r *UserRepository) DeleteJoin(ctx context.Context, id uint, channel_id uint) error {
	var userJoin gormmodel.MiddleChannelUser
	if err := r.gorm.WithContext(ctx).Where("user_id = ? && channel_id = ?", id, channel_id).Delete(&userJoin).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_CHANNEL_TABLE + fmt.Sprintf("%d:%d", id, channel_id)
	go func(ctx context.Context, key string) {
		if err := r.redis.Del(ctx, key).Err(); err != nil {
			log.Printf("redis Del error: %v", err)
		}
	}(ctx, key)
	return nil
}

func (r *UserRepository) CreateParticate(ctx context.Context, id uint, event_id uint, userParticate entities.UserParticate) error {
	var user gormmodel.User
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}
	var event gormmodel.Event
	if err := r.gorm.WithContext(ctx).Where("id = ?", event_id).First(&event).Error; err != nil {
		return err
	}

	userJoinModel := gormmodel.MiddleEventUser{
		UserID:        id,
		EventID:       event_id,
		Role:          userParticate.Role,
		LastParticate: userParticate.LastParticate,
	}

	if err := r.gorm.WithContext(ctx).Create(&userJoinModel).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_EVENT_TABLE + fmt.Sprintf("%d:%d", id, event_id)
	fields := map[string]interface{}{
		rediskey.REDIS_USER_EVENT_FIELD_LASTPARTICATE: userParticate.LastParticate,
		rediskey.REDIS_USER_EVENT_FIELD_ROLE:          userParticate.Role,
	}

	go func(ctx context.Context, key string, fields map[string]interface{}) {
		if err := r.redis.HSet(ctx, key, fields).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 24*time.Hour)
	}(ctx, key, fields)
	return nil
}

func (r *UserRepository) UpdateParticateTime(ctx context.Context, id uint, event_id uint, particateTime int64) (*entities.UserParticate, error) {
	var userEvent gormmodel.MiddleEventUser
	if err := r.gorm.WithContext(ctx).Where("user_id = ? && event_id =?", id, event_id).First(&userEvent).Update("last_particate", particateTime).Error; err != nil {
		return nil, err
	}
	key := rediskey.REDIS_USER_EVENT_TABLE + fmt.Sprintf("%d:%d", id, event_id)
	go func(ctx context.Context, key string, lastParticate int64) {
		if err := r.redis.HSet(ctx, key, rediskey.REDIS_USER_EVENT_FIELD_LASTPARTICATE, lastParticate).Err(); err != nil {
			log.Printf("redis HSet error: %v", err)
		}
		r.redis.Expire(ctx, key, 7*24*time.Hour)
	}(ctx, key, particateTime)
	return userEvent.ToDomain(), nil
}

func (r *UserRepository) DeleteParticate(ctx context.Context, id uint, event_id uint) error {
	var userParticate gormmodel.MiddleEventUser
	if err := r.gorm.WithContext(ctx).Where("user_id = ? && event_id = ?", id, event_id).Delete(&userParticate).Error; err != nil {
		return err
	}
	key := rediskey.REDIS_USER_EVENT_TABLE + fmt.Sprintf("%d:%d", id, event_id)
	go func(ctx context.Context, key string) {
		if err := r.redis.Del(ctx, key).Err(); err != nil {
			log.Printf("redis Del error: %v", err)
		}
	}(ctx, key)
	return nil
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
