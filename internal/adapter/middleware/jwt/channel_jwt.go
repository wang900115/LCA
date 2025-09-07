package jwt

import (
	"errors"
	"strings"

	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/implement"

	"github.com/gin-gonic/gin"
)

type CHANNELJWT struct {
	response iresponse.IResponse
	token    implement.TokenImplement
}

func NewCHANNELJWT(response iresponse.IResponse, token implement.TokenImplement) *CHANNELJWT {
	return &CHANNELJWT{response: response, token: token}
}

func (j *CHANNELJWT) Middleware(c *gin.Context) {
	token, err := j.extract(c)
	if err != nil {
		j.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	tokenClaims, err := j.token.ValidateChannelToken(token)
	if err != nil {
		j.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	if tokenClaims.UserID != c.GetUint("user_id") {
		c.Abort()
		return
	}

	c.Set("channel_id", tokenClaims.ChannelID)
	c.Set("role", tokenClaims.JoinStatus.Role)
	c.Set("last_join", tokenClaims.JoinStatus.LastJoin)
	c.Next()

}

func (CHANNELJWT) extract(c *gin.Context) (string, error) {
	authorization := c.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("no token")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid authorization format")
	}

	return strings.TrimPrefix(authorization, "Bearer "), nil
}
