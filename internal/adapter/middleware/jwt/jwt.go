package jwt

import (
	"errors"
	"strings"

	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type USERJWT struct {
	response iresponse.IResponse
	token    usecase.TokenUsecase
}

func NewJWT(response iresponse.IResponse, token *usecase.TokenUsecase) *USERJWT {
	return &USERJWT{response: response, token: *token}
}

func (j *USERJWT) Middleware(c *gin.Context) {
	token, err := j.extract(c)
	if err != nil {
		j.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	tokenClaims, err := j.token.ValidateUserToken(token)
	if err != nil {
		j.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	c.Set("user_id", tokenClaims.UserID)
	c.Set("ip_address", tokenClaims.LoginStatus.IPAddress)
	c.Set("expired_at", tokenClaims.ExpiredAt)

	c.Next()

}

func (USERJWT) extract(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("no token")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization format")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}
