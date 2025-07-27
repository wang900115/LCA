package redismodel

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/wang900115/LCA/pkg/domain"
)

/*
當 access token 過期時：

查找 refresh-token:<user_id>:<token_id> → 取得 salt

確認跟 access token 的 salt 是否一致

若一致 → 重新簽發 access + refresh token

若不一致 → 有潛在風險（可能是被替換或駭入）
*/

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (a AccessTokenClaims) ToDomain() domain.TokenClaims {
	return domain.TokenClaims{
		Name:             "access",
		UserID:           a.UserID,
		Role:             a.Role,
		RegisteredClaims: a.RegisteredClaims,
	}
}

func (a AccessTokenClaims) FromDomain(token domain.TokenClaims) AccessTokenClaims {
	return AccessTokenClaims{
		UserID:           token.UserID,
		Role:             token.Role,
		RegisteredClaims: token.RegisteredClaims,
	}
}

func (r RefreshTokenClaims) ToDomain() domain.TokenClaims {
	return domain.TokenClaims{
		Name:             "refresh",
		UserID:           r.UserID,
		RegisteredClaims: r.RegisteredClaims,
	}
}

func (r RefreshTokenClaims) FromDomain(token domain.TokenClaims) RefreshTokenClaims {
	return RefreshTokenClaims{
		UserID:           token.UserID,
		RegisteredClaims: token.RegisteredClaims,
	}
}

// 驗簽
// func (a *AccessTokenClaims) ParseAccessToken(tokenString string, fullSecret string) (*AccessTokenClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
// 		return fullSecret, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	claims, ok := token.Claims.(*AccessTokenClaims)
// 	if !ok || !token.Valid {
// 		return nil, common.Unauthorized
// 	}
// 	return claims, nil
// }

// 驗簽
// func (r *RefreshTokenClaims) ParseRefreshToken(tokenString string, fullSecret string) (*RefreshTokenClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
// 		return fullSecret, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	claims, ok := token.Claims.(*RefreshTokenClaims)
// 	if !ok || token.Valid {
// 		return nil, common.Unauthorized
// 	}
// 	return claims, nil
// }

// func buildFullSecret(secret string) func(string) string {
// 	return func(salt string) string {
// 		h := hmac.New(sha256.New, []byte(secret))
// 		h.Write([]byte(salt))
// 		return hex.EncodeToString(h.Sum(nil))
// 	}
// }

// func groupAccessToken(userID, tokenID string) string {
// 	return rediskey.REDIS_TABLE_ACCESS_TOKEN + userID + ":" + tokenID
// }

// func groupRefreshToken(userID, tokenID string) string {
// 	return rediskey.REDIS_TABLE_REFRESH_TOKEN + userID + ":" + tokenID
// }
