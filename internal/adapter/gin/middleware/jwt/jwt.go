package jwt

import (
	"errors"
	"strings"

	iresponse "github.com/wang900115/LCA/internal/adapter/gin/controller/response"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type JWT struct {
	response iresponse.IResponse
	token    usecase.TokenUsecase
}

func NewJWT(response iresponse.IResponse, token *usecase.TokenUsecase) *JWT {
	return &JWT{response: response, token: *token}
}

func (j *JWT) Middleware(c *gin.Context) {
	token, err := j.extract(c)
	if err != nil {
		j.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	tokenClaims, err := j.token.ValidateToken(token)
	if err != nil {
		j.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	c.Set("channel", tokenClaims.Channel)
	c.Set("user", tokenClaims.User)
	c.Set("expired_at", tokenClaims.ExpiredAt)

	c.Next()

}

func (JWT) extract(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("no token")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization format")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}
