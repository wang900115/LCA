package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/wang900115/LCA/internal/adapter/gorm/model"
	"github.com/wang900115/LCA/internal/domain/entities"
	"github.com/wang900115/LCA/internal/domain/irepository"

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

type TokenRepository struct {
	redis      *redis.Client
	expiration time.Duration
}

func NewTokenRepository(redis *redis.Client, expiration time.Duration) irepository.ITokenRepository {
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
	tokenClaimsModel := model.TokenClaims{
		UserUUID:    tokenClaims.UserUUID,
		ChannelUUID: tokenClaims.ChannelUUID,
		Username:    tokenClaims.Username,
	}
	tokenClaimsModel.ExpiresAt = jwt.NewNumericDate(time.Now().Add(r.expiration))

	_, err := r.redis.Set(context.Background(), jwtsaltPrefix+tokenClaims.UserUUID+tokenClaims.ChannelUUID, string(salt), r.expiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaimsModel).SignedString(append([]byte(tokenClaimsModel.UserUUID+tokenClaims.ChannelUUID), salt...))
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

	userUUID, ok := mapClaims["user_uuid"].(string)
	if !ok {
		return entities.TokenClaims{}, errors.New("token map userUUID failed")
	}

	channelUUID, ok := mapClaims["channel_uuid"].(string)
	if !ok {
		return entities.TokenClaims{}, errors.New("token map channelUUID failed")
	}

	salt, err := r.redis.Get(context.Background(), jwtsaltPrefix+userUUID+channelUUID).Result()
	if err != nil {
		return entities.TokenClaims{}, err
	}

	key := []byte(userUUID + channelUUID + salt)
	tokenClaims, parseErr := jwt.ParseWithClaims(token, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if parseErr != nil {
		return entities.TokenClaims{}, parseErr
	}

	if !tokenClaims.Valid {
		return entities.TokenClaims{}, ErrTokenExpired
	}

	tokenClaimsModel, ok := tokenClaims.Claims.(*model.TokenClaims)
	if !ok {
		return entities.TokenClaims{}, ErrInvalidToken
	}
	return tokenClaimsModel.ToDomain(), nil
}

func (r *TokenRepository) DeleteToken(userUUID, channelUUID string) error {
	return r.redis.Del(context.Background(), jwtsaltPrefix+userUUID+channelUUID).Err()

}
