package middleware

import (
	"errors"
	"strings"

	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
	"github.com/wang900115/LCA/pkg/implement"

	"github.com/gin-gonic/gin"
)

type JWT struct {
	resp   iresponse.IResponse
	token  implement.TokenAuthService
	secret string
}

func NewJWT(resp iresponse.IResponse, token *implement.TokenAuthService, secret string) *JWT {
	return &JWT{resp: resp, token: *token, secret: secret}
}

func (j *JWT) Middleware(c *gin.Context) {
	accessToken, err := j.extractAccessToken(c)
	if err != nil {
		j.resp.FailWithError(c, common.PARAM_ERROR, err)
		c.Abort()
		return
	}

	tokenClaims, err := j.token.ValidateAccess(c, accessToken, j.secret)
	if err != nil {
		if errors.Is(err, common.TokenExpired) {
			refreshToken, err := j.extractRefreshToken(c)
			if err != nil {
				j.resp.FailWithError(c, common.UNAUTHORIZED_FAIL, err)
				c.Abort()
				return
			}

			tokenClaims, err = j.token.ValidateRefresh(c, refreshToken, j.secret)
			if err != nil {
				j.resp.FailWithError(c, common.UNAUTHORIZED_FAIL, err)
				c.Abort()
				return
			}

			newAccessToken, err := j.token.Refresh(c, tokenClaims.UserID, j.secret)
			if err != nil {
				j.resp.FailWithError(c, common.UNAUTHORIZED_FAIL, err)
				c.Abort()
				return
			}
			c.Header("X-Access-Token", newAccessToken)
		} else {
			j.resp.FailWithError(c, common.UNAUTHORIZED_FAIL, err)
			c.Abort()
			return
		}
	}

	c.Set("user_id", tokenClaims.UserID)
	c.Next()
}

func (JWT) extractAccessToken(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("no token")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization format")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}

func (JWT) extractRefreshToken(c *gin.Context) (string, error) {
	token := c.GetHeader("X-Refresh-Token")
	if token == "" {
		return "", common.TokenMissed
	}
	return token, nil
}
