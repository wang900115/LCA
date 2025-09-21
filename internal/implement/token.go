package implement

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"strconv"
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
	CreateEventToken(context.Context, entities.EventTokenClaims) (string, error)
	ValidateUserToken(string) (*entities.UserTokenClaims, error)
	ValidateChannelToken(string) (*entities.ChannelTokenClaims, error)
	ValidateEventToken(string) (*entities.EventTokenClaims, error)
	DeleteUserToken(context.Context, uint) error
	DeleteChannelToken(context.Context, uint, uint) error
	DeleteEventToken(context.Context, uint, uint) error
}

type TokenRepository struct {
	redis               *redis.Client
	loginExpiration     time.Duration
	joinExpiration      time.Duration
	particateExpiration time.Duration
	loginSecret         []byte
	joinSecret          []byte
	particateSecret     []byte
}

func NewTokenRepository(redis *redis.Client, loginExpiration time.Duration, joinExpiration time.Duration, particateExpiration time.Duration, loginSecret []byte, joinSecret []byte, particateSecret []byte) TokenImplement {
	return &TokenRepository{
		redis:               redis,
		loginExpiration:     loginExpiration,
		joinExpiration:      joinExpiration,
		particateExpiration: particateExpiration,
		loginSecret:         loginSecret,
		joinSecret:          joinSecret,
		particateSecret:     particateSecret,
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
	id := strconv.FormatUint(uint64(tokenClaims.UserID), 10)
	_, err := r.redis.Set(ctx, jwtsaltPrefix+id, string(salt), r.loginExpiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaimsModel).SignedString(append(r.loginSecret, salt...))
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
	id := strconv.FormatUint(uint64(channelClaims.UserID), 10)
	channelId := strconv.FormatUint(uint64(channelClaims.ChannelID), 10)
	_, err := r.redis.Set(ctx, jwtsaltPrefix+id+channelId, string(salt), r.joinExpiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, channelClaimsModel).SignedString(append(r.joinSecret, salt...))
}

func (r *TokenRepository) CreateEventToken(ctx context.Context, eventClaims entities.EventTokenClaims) (string, error) {
	salt := r.GenerateSalt(saltSize)
	eventClaimsModel := redismodel.EventTokenClaims{
		UserID:        eventClaims.UserID,
		EventID:       eventClaims.EventID,
		Role:          eventClaims.ParticateStatus.Role,
		LastParticate: eventClaims.ParticateStatus.LastParticate,
	}
	eventClaimsModel.ExpiresAt = jwt.NewNumericDate(time.Now().Add(r.joinExpiration))
	id := strconv.FormatUint(uint64(eventClaims.UserID), 10)
	eventId := strconv.FormatUint(uint64(eventClaims.EventID), 10)
	_, err := r.redis.Set(ctx, jwtsaltPrefix+id+eventId, string(salt), r.joinExpiration).Result()
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, eventClaimsModel).SignedString(append(r.particateSecret, salt...))
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

func (r *TokenRepository) ValidateEventToken(token string) (*entities.EventTokenClaims, error) {
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
	event, ok := mapClaims["event"].(string)
	if !ok {
		return nil, errors.New("token map event failed")
	}

	salt, err := r.redis.Get(context.Background(), jwtsaltPrefix+user+event).Result()
	if err != nil {
		return nil, err
	}
	key := bytes.Join([][]byte{r.particateSecret, []byte(salt)}, []byte{})
	tokenClaims, parseErr := jwt.ParseWithClaims(token, &redismodel.ChannelTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if parseErr != nil {
		return nil, parseErr
	}
	if !tokenClaims.Valid {
		return nil, ErrTokenExpired
	}
	tokenClaimsModel, ok := tokenClaims.Claims.(*redismodel.EventTokenClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return tokenClaimsModel.ToDomain(), nil
}

func (r *TokenRepository) DeleteUserToken(ctx context.Context, userId uint) error {
	id := strconv.FormatUint(uint64(userId), 10)
	return r.redis.Del(context.Background(), jwtsaltPrefix+id).Err()
}

func (r *TokenRepository) DeleteChannelToken(ctx context.Context, userId uint, channelId uint) error {
	id := strconv.FormatUint(uint64(userId), 10)
	channel_id := strconv.FormatUint(uint64(channelId), 10)
	return r.redis.Del(context.Background(), jwtsaltPrefix+id+channel_id).Err()
}

func (r *TokenRepository) DeleteEventToken(ctx context.Context, userId uint, eventId uint) error {
	id := strconv.FormatUint(uint64(userId), 10)
	event_id := strconv.FormatUint(uint64(eventId), 10)
	return r.redis.Del(context.Background(), jwtsaltPrefix+id+event_id).Err()
}
