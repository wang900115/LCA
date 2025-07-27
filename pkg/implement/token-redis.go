package implement

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wang900115/LCA/pkg/common"
	"github.com/wang900115/LCA/pkg/domain"
	rediskey "github.com/wang900115/LCA/pkg/redis/key"
	redismodel "github.com/wang900115/LCA/pkg/redis/model"
	"go.uber.org/zap"
)

type TokenAuthService interface {
	// 簽發
	Generate(c context.Context, userID string, role string, secretKey string) (accessToken string, refreshToken string, err error)

	// 驗證ACCESS
	ValidateAccess(c context.Context, tokenString string, secretKey string) (domain.TokenClaims, error)

	// 驗證REFRESH
	ValidateRefresh(c context.Context, tokenString string, secretKey string) (domain.TokenClaims, error)

	// 重發
	Refresh(c context.Context, userID string, secretKey string) (newAccessToken string, err error)

	// 刪除
	Delete(c context.Context, userID string) error
}

type TokenAuthRepository struct {
	redis  *redis.Client
	logger *zap.Logger
}

func NewTokenAuthRepository(redis *redis.Client, logger *zap.Logger) TokenAuthService {
	return &TokenAuthRepository{redis: redis, logger: logger}
}

func (t *TokenAuthRepository) Generate(c context.Context, userID string, role string, secretKey string) (accessToken string, refreshToken string, err error) {

	now := time.Now().UTC()
	accessTokenExp := jwt.NewNumericDate(now.Add(24 * time.Hour))
	refreshTokenExp := jwt.NewNumericDate(now.Add(7 * 24 * time.Hour))
	nowAt := jwt.NewNumericDate(now)

	salt := generateSalt(now, secretKey)

	accessClaims := redismodel.AccessTokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: accessTokenExp,
			IssuedAt:  nowAt,
			NotBefore: nowAt,
			ID:        uuid.NewString(),
		},
	}

	refreshClaims := redismodel.RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: refreshTokenExp,
			IssuedAt:  nowAt,
			NotBefore: nowAt,
			ID:        uuid.NewString(),
		},
	}

	if err := t.redis.Set(c, rediskey.REDIS_STRING_ACCESS_TOKEN+accessClaims.UserID+accessClaims.ID, salt, 24*time.Hour).Err(); err != nil {
		t.logger.Error("Redis Write Access Token String Err", zap.Error(err))
		return "", "", err
	}

	if err := t.redis.Set(c, rediskey.REDIS_STRING_REFRESH_TOKEN+refreshClaims.UserID+refreshClaims.ID, salt, 7*24*time.Hour).Err(); err != nil {
		t.logger.Error("Redis Write Refresh Token String Err", zap.Error(err))
		return "", "", err
	}

	accessTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signature := groupName(salt, secretKey)
	signedAccess, err := accessTokenString.SignedString([]byte(signature))
	if err != nil {
		return "", "", err
	}
	signedReFresh, err := refreshTokenString.SignedString([]byte(signature))
	if err != nil {
		return "", "", err
	}

	return signedAccess, signedReFresh, nil

}

func (t *TokenAuthRepository) ValidateAccess(c context.Context, tokenString string, secretKey string) (domain.TokenClaims, error) {
	claims := &redismodel.AccessTokenClaims{}
	_, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(tokenString, claims)
	if err != nil {
		return domain.TokenClaims{}, err
	}
	salt, err := t.redis.Get(c, rediskey.REDIS_STRING_ACCESS_TOKEN+claims.UserID+claims.ID).Result()
	if err != nil {
		return domain.TokenClaims{}, err
	}

	signature := groupName(salt, secretKey)
	claims = &redismodel.AccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})
	if err != nil {
		if errors.Is(err, common.TokenExpired) {
			return domain.TokenClaims{}, common.TokenExpired
		}
		return domain.TokenClaims{}, common.TokenExpired
	}

	if !token.Valid {
		return domain.TokenClaims{}, common.TokenInvalid
	}

	return claims.ToDomain(), nil
}

func (t *TokenAuthRepository) ValidateRefresh(c context.Context, tokenString string, secretKey string) (domain.TokenClaims, error) {
	claims := &redismodel.RefreshTokenClaims{}
	_, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(tokenString, claims)
	if err != nil {
		return domain.TokenClaims{}, err
	}
	salt, err := t.redis.Get(c, rediskey.REDIS_STRING_ACCESS_TOKEN+claims.UserID+claims.ID).Result()
	if err != nil {
		return domain.TokenClaims{}, err
	}

	signature := groupName(salt, secretKey)
	claims = &redismodel.RefreshTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})
	if err != nil {
		if errors.Is(err, common.TokenExpired) {
			return domain.TokenClaims{}, common.TokenExpired
		}
		return domain.TokenClaims{}, common.TokenExpired
	}

	if !token.Valid {
		return domain.TokenClaims{}, common.TokenInvalid
	}
	return claims.ToDomain(), nil
}

// 不需要更新Refresh
func (t *TokenAuthRepository) Refresh(c context.Context, userID string, secretKey string) (newAccessToken string, err error) {
	now := time.Now().UTC()
	accessTokenExp := jwt.NewNumericDate(now.Add(24 * time.Hour))
	nowAt := jwt.NewNumericDate(now)

	salt := generateSalt(now, secretKey)

	accessClaims := redismodel.AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: accessTokenExp,
			IssuedAt:  nowAt,
			NotBefore: nowAt,
			ID:        uuid.NewString(),
		},
	}

	if err := t.redis.Set(c, rediskey.REDIS_STRING_ACCESS_TOKEN+accessClaims.UserID+accessClaims.ID, salt, 24*time.Hour).Err(); err != nil {
		t.logger.Error("Redis Write Refresh Access Token String Err", zap.Error(err))
		return "", err
	}

	accessTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signature := groupName(salt, secretKey)
	signedAccess, err := accessTokenString.SignedString([]byte(signature))
	if err != nil {
		return "", err
	}
	return signedAccess, nil
}

func (t *TokenAuthRepository) Delete(c context.Context, userID string) error {
	accessPattern := rediskey.REDIS_STRING_ACCESS_TOKEN + userID + "*"
	refreshPattern := rediskey.REDIS_STRING_REFRESH_TOKEN + userID + "*"

	iterAccess := t.redis.Scan(c, 0, accessPattern, 0).Iterator()
	for iterAccess.Next(c) {
		key := iterAccess.Val()
		if err := t.redis.Del(c, key).Err(); err != nil {
			t.logger.Error("failed to delete access token key from redis", zap.String("key", key), zap.Error(err))
			return err
		}
	}

	if err := iterAccess.Err(); err != nil {
		t.logger.Error("error during redis scan access-token", zap.Error(err))
		return err
	}

	iterRefresh := t.redis.Scan(c, 0, refreshPattern, 0).Iterator()
	for iterRefresh.Next(c) {
		key := iterAccess.Val()
		if err := t.redis.Del(c, key).Err(); err != nil {
			t.logger.Error("failed to delete refresh token key from redis", zap.String("key", key), zap.Error(err))
			return err
		}
	}

	if err := iterRefresh.Err(); err != nil {
		t.logger.Error("error during redis scan refresh-token", zap.Error(err))
		return err
	}

	return nil
}

func generateSalt(t time.Time, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(t.Format(time.RFC3339Nano)))
	return hex.EncodeToString(h.Sum(nil))
}

func groupName(salt, secretKey string) string {
	return salt + secretKey
}
