package router

import (
	"github.com/wang900115/LCA/internal/adapter/controller"

	"github.com/gin-gonic/gin"
)

type ChannelRouter struct {
	channelController controller.ChannelController
}

func NewChannelRouter(channelController *controller.ChannelController) IRoute {
	return &ChannelRouter{channelController: *channelController}
}

func (cr *ChannelRouter) Setup(router *gin.RouterGroup) {
	channelGroup := router.Group("v1/channel/")
	{
		channelGroup.POST("/create", cr.channelController.CreateChannel)
		channelGroup.GET("/queryuser", cr.channelController.QueryUsers)
		channelGroup.GET("/query", cr.channelController.QueryChannel)
	}
}
