package jwt

import (
	"errors"
	"strings"

	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/implement"

	"github.com/gin-gonic/gin"
)

type USERJWT struct {
	response iresponse.IResponse
	token    implement.TokenImplement
}

func NewUSERJWT(response iresponse.IResponse, token implement.TokenImplement) *USERJWT {
	return &USERJWT{response: response, token: token}
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
	c.Set("last_login", tokenClaims.LoginStatus.LastLogin)

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
