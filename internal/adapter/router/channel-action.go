package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/controller"
	"github.com/wang900115/LCA/internal/adapter/middleware/jwt"
)

type ChannelRouter struct {
	channelController controller.ChannelController
	userJWT           jwt.USERJWT
	channelJWT        jwt.CHANNELJWT
}

func NewChannelRouter(channelController *controller.ChannelController, userJWT *jwt.USERJWT, CHANNELJWT *jwt.CHANNELJWT) IRoute {
	return &ChannelRouter{channelController: *channelController, userJWT: *userJWT, channelJWT: *CHANNELJWT}
}

func (cr *ChannelRouter) Setup(router *gin.RouterGroup) {
	channelGroup := router.Group("v1/channel/", cr.userJWT.Middleware)
	{
		channelGroup.POST("create", cr.channelController.CreateChannel)
	}

	channelInfoGroup := router.Group("v1/channel/info", cr.userJWT.Middleware, cr.channelJWT.Middleware)
	{
		channelInfoGroup.POST("users", cr.channelController.QueryUsers)
		channelInfoGroup.POST("messages", cr.channelController.QueryMessages)
	}
}
