package route

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/controller"
	"github.com/wang900115/LCA/pkg/common/router"
)

type ChannelMessageRouter struct {
	channelMessage controller.ChannelMessageController
}

func NewChannelMessageRouter(channelMessage *controller.ChannelMessageController) router.IRoute {
	return &ChannelMessageRouter{channelMessage: *channelMessage}
}

func (c *ChannelMessageRouter) Setup(router *gin.RouterGroup) {
	channelMessageGroup := router.Group("v1/channel-message/")
	{
		channelMessageGroup.POST("/query", c.channelMessage.GetChannelMessages) // 獲取頻道的訊息
	}
}
