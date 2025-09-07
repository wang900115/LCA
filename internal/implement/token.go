package implement

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"time"

	redismodel "github.com/wang900115/LCA/internal/adapter/redis/model"
	"github.com/wang900115/LCA/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

const (
	jwtsaltPrefix = "jwtsalt:"
	saltSize      = 16
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type TokenImplement interface {
	GenerateSalt(int) []byte
	CreateUserToken(context.Context, entities.UserTokenClaims) (string, error)
	CreateChannelToken(context.Context, entities.ChannelTokenClaims) (string, error)
	ValidateUserToken(string) (*entities.UserTokenClaims, error)
	ValidateChannelToken(string) (*entities.ChannelTokenClaims, error)
	DeleteUserToken(context.Context, uint) error
	DeleteChannelToken(context.Context, uint, uint) error
}

type TokenRepository struct {
	redis           *redis.Client
	loginExpiration time.Duration
	joinExpiration  time.Duration
	loginSecret     []byte
	joinSecret      []byte
}

func NewTokenRepository(redis *redis.Client, loginExpiration time.Duration, joinExpiration time.Duration, loginSecret []byte, joinSecret []byte) TokenImplement {
	return &TokenRepository{
		redis:           redis,
		loginExpiration: loginExpiration,
		joinExpiration:  joinExpiration,
		loginSecret:     loginSecret,
		joinSecret:      joinSecret,
	}
}

func (TokenRepository) GenerateSalt(saltSize int) []byte {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	return salt
}

func (r *TokenRepository) CreateUserToken(ctx context.Context, tokenClaims entities.UserTokenClaims) (string, error) {
	salt := r.GenerateSalt(saltSize)
	tokenClaimsModel := redismodel.UserTokenClaims{
		UserID:     tokenClaims.UserID,
		LastLogin:  tokenClaims.LoginStatus.LastLogin,
		IPAddress:  *tokenClaims.LoginStatus.IPAddress,
		DeviceInfo: *tokenClaims.LoginStatus.DeviceInfo,
	}
	tokenClaimsModel.ExpiresAt = jwt.NewNumericDate(time.Now().Add(r.loginExpiration))

	_, err := r.redis.Set(ctx, jwtsaltPrefix+string(tokenClaims.UserID), string(salt), r.loginExpiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaimsModel).SignedString(append([]byte(string(tokenClaimsModel.UserID)), salt...))
}

func (r *TokenRepository) CreateChannelToken(ctx context.Context, channelClaims entities.ChannelTokenClaims) (string, error) {
	salt := r.GenerateSalt(saltSize)
	channelClaimsModel := redismodel.ChannelTokenClaims{
		UserID:    channelClaims.UserID,
		ChannelID: channelClaims.ChannelID,
		Role:      channelClaims.JoinStatus.Role,
		LastJoin:  channelClaims.JoinStatus.LastJoin,
	}
	channelClaimsModel.ExpiresAt = jwt.NewNumericDate(time.Now().Add(r.joinExpiration))

	_, err := r.redis.Set(ctx, jwtsaltPrefix+string(channelClaims.UserID)+string(channelClaims.ChannelID), string(salt), r.joinExpiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, channelClaimsModel).SignedString(append(r.joinSecret, salt...))
}

func (r *TokenRepository) ValidateUserToken(token string) (*entities.UserTokenClaims, error) {
	unvertifiedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	mapClaims, ok := unvertifiedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token map failed")
	}
	user, ok := mapClaims["user"].(string)
	if !ok {
		return nil, errors.New("token map user failed")
	}
	salt, err := r.redis.Get(context.Background(), jwtsaltPrefix+user).Result()
	if err != nil {
		return nil, err
	}
	key := bytes.Join([][]byte{r.loginSecret, []byte(salt)}, []byte{})
	tokenClaims, parseErr := jwt.ParseWithClaims(token, &redismodel.UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if !tokenClaims.Valid {
		return nil, ErrTokenExpired
	}
	tokenClaimsModel, ok := tokenClaims.Claims.(*redismodel.UserTokenClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return tokenClaimsModel.ToDomain(), nil
}

func (r *TokenRepository) ValidateChannelToken(token string) (*entities.ChannelTokenClaims, error) {
	unvertifiedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	mapClaims, ok := unvertifiedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token map failed")
	}
	user, ok := mapClaims["user"].(string)
	if !ok {
		return nil, errors.New("token map user failed")
	}
	channel, ok := mapClaims["channel"].(string)
	if !ok {
		return nil, errors.New("token map channel failed")
	}

	salt, err := r.redis.Get(context.Background(), jwtsaltPrefix+user+channel).Result()
	if err != nil {
		return nil, err
	}
	key := bytes.Join([][]byte{r.joinSecret, []byte(salt)}, []byte{})
	tokenClaims, parseErr := jwt.ParseWithClaims(token, &redismodel.ChannelTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if !tokenClaims.Valid {
		return nil, ErrTokenExpired
	}
	tokenClaimsModel, ok := tokenClaims.Claims.(*redismodel.ChannelTokenClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return tokenClaimsModel.ToDomain(), nil
}

func (r *TokenRepository) DeleteUserToken(ctx context.Context, userId uint) error {
	return r.redis.Del(context.Background(), jwtsaltPrefix+string(userId)).Err()
}

func (r *TokenRepository) DeleteChannelToken(ctx context.Context, userId uint, channelId uint) error {
	return r.redis.Del(context.Background(), jwtsaltPrefix+string(userId)+string(channelId)).Err()
}
