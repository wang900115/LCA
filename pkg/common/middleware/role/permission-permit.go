package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
	"github.com/wang900115/LCA/pkg/implement"
)

type Permission struct {
	resp   iresponse.IResponse
	token  implement.TokenAuthService
	secret string
}

func NewPermission(resp iresponse.IResponse, token *implement.TokenAuthService, secret string) *Permission {
	return &Permission{resp: resp, token: *token, secret: secret}
}

func (p *Permission) RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := p.extractAccessToken(c)
		if err != nil {
			p.resp.FailWithError(c, common.PARAM_ERROR, err)
			c.Abort()
			return
		}

		tokenClaims, err := p.token.ValidateAccess(c, accessToken, p.secret)
		if err != nil {
			p.resp.FailWithError(c, common.UNAUTHORIZED_FAIL, err)
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if tokenClaims.Role == role {
				c.Next()
				return
			}
		}
		p.resp.FailWithError(c, common.PERMISSION_ACCESS, errors.New("insufficient permission"))
		c.Abort()
	}
}

func (Permission) extractAccessToken(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("no token")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization format")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}
