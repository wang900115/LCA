package router

import (
	"LCA/internal/adapter/gin/controller"

	"github.com/gin-gonic/gin"
)

type ChannelRouter struct {
	channelController controller.ChannelController
}

func NewChannelRouter(channelController controller.ChannelController) *ChannelRouter {
	return &ChannelRouter{channelController: channelController}
}

func (cr *ChannelRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("v1/Channel/")
	{
		userGroup.POST("/create", cr.channelController.CreateChannel)
		userGroup.POST("/queryuser", cr.channelController.QueryUsers)
		userGroup.POST("/query", cr.channelController.QueryChannel)
	}
}
