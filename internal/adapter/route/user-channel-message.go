package route

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/controller"
	middleware "github.com/wang900115/LCA/pkg/common/middleware/jwt"
	"github.com/wang900115/LCA/pkg/common/router"
)

type UserChannelMessageRouter struct {
	userMessage controller.UserChannelMessageController
	jwt         middleware.JWT
}

func NewUserChannelMessageRouter(userMessage *controller.UserChannelMessageController, jwt *middleware.JWT) router.IRoute {
	return &UserChannelMessageRouter{userMessage: *userMessage, jwt: *jwt}
}

func (u *UserChannelMessageRouter) Setup(router *gin.RouterGroup) {
	userChannelMessageGroup := router.Group("v1/user-channel-message/", u.jwt.Middleware)
	{
		userChannelMessageGroup.POST("/query", u.userMessage.GetChannelUserMessages) // 用戶訊息
	}
}
