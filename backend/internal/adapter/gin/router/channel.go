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
	channelGroup := router.Group("v1/Channel/")
	{
		channelGroup.POST("/create", cr.channelController.CreateChannel)
		channelGroup.GET("/query", cr.channelController.DeleteChannel)
		channelGroup.DELETE("/delete", cr.channelController.QueryChannel)
	}
}
