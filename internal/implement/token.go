package implement

import (
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
	CreateUserToken(entities.UserTokenClaims) (string, error)
	CreateChannelToken(entities.ChannelTokenClaims) (string, error)
	ValidateUserToken(string) (entities.UserTokenClaims, error)
	ValidateChannelToken(string) (entities.ChannelTokenClaims, error)

	DeleteUserToken(uint) error
	DeleteChannelToken(uint, uint) error
}
type TokenRepository struct {
	redis      *redis.Client
	expiration time.Duration
}

func NewTokenRepository(redis *redis.Client, expiration time.Duration) TokenImplement {
	return &TokenRepository{
		redis:      redis,
		expiration: expiration,
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

func (r *TokenRepository) CreateToken(tokenClaims entities.TokenClaims) (string, error) {
	salt := r.GenerateSalt(saltSize)
	tokenClaimsModel := redismodel.TokenClaims{
		User:    tokenClaims.User,
		Channel: tokenClaims.Channel,
	}
	tokenClaimsModel.ExpiresAt = jwt.NewNumericDate(time.Now().Add(r.expiration))

	_, err := r.redis.Set(context.Background(), jwtsaltPrefix+tokenClaims.User+tokenClaims.Channel, string(salt), r.expiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaimsModel).SignedString(append([]byte(tokenClaimsModel.User+tokenClaims.Channel), salt...))
}

func (r *TokenRepository) ValidateToken(token string) (entities.TokenClaims, error) {
	unvertifiedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return entities.TokenClaims{}, err
	}
	mapClaims, ok := unvertifiedToken.Claims.(jwt.MapClaims)
	if !ok {
		return entities.TokenClaims{}, errors.New("token map failed")
	}

	user, ok := mapClaims["user"].(string)
	if !ok {
		return entities.TokenClaims{}, errors.New("token map userUUID failed")
	}

	channel, ok := mapClaims["channel"].(string)
	if !ok {
		return entities.TokenClaims{}, errors.New("token map channelUUID failed")
	}

	salt, err := r.redis.Get(context.Background(), jwtsaltPrefix+user+channel).Result()
	if err != nil {
		return entities.TokenClaims{}, err
	}

	key := []byte(user + channel + salt)
	tokenClaims, parseErr := jwt.ParseWithClaims(token, &redismodel.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if parseErr != nil {
		return entities.TokenClaims{}, parseErr
	}

	if !tokenClaims.Valid {
		return entities.TokenClaims{}, ErrTokenExpired
	}

	tokenClaimsModel, ok := tokenClaims.Claims.(*redismodel.TokenClaims)
	if !ok {
		return entities.TokenClaims{}, ErrInvalidToken
	}
	return tokenClaimsModel.ToDomain(), nil
}

func (r *TokenRepository) DeleteUserToken(user string) error {
	return r.redis.Del(context.Background(), jwtsaltPrefix+user).Err()
}
