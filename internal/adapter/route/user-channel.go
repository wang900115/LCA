package route

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/controller"
	middleware "github.com/wang900115/LCA/pkg/common/middleware/jwt"
	"github.com/wang900115/LCA/pkg/common/router"
)

type UserChannelRouter struct {
	userChannel controller.UserChannelController
	jwt         middleware.JWT
}

func NewUserChannelRouter(userChannel *controller.UserChannelController, jwt *middleware.JWT) router.IRoute {
	return &UserChannelRouter{userChannel: *userChannel, jwt: *jwt}
}

func (u *UserChannelRouter) Setup(router *gin.RouterGroup) {
	userChannelGroup := router.Group("v1/user-channel/", u.jwt.Middleware)
	{
		userChannelGroup.POST("/query", u.userChannel.GetUserChannels) // 獲取該user的所有頻道
	}
}
