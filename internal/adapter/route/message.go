package route

import (
	"github.com/wang900115/LCA/internal/adapter/controller"
	middleware "github.com/wang900115/LCA/pkg/common/middleware/jwt"
	"github.com/wang900115/LCA/pkg/common/router"

	"github.com/gin-gonic/gin"
)

type MessageRouter struct {
	message controller.MessageController
	jwt     middleware.JWT
}

func NewMessageRouter(message *controller.MessageController, jwt *middleware.JWT) router.IRoute {
	return &MessageRouter{message: *message, jwt: *jwt}
}

func (mr *MessageRouter) Setup(router *gin.RouterGroup) {
	messageGroup := router.Group("v1/message/", mr.jwt.Middleware)
	{
		messageGroup.POST("/create", mr.message.Create) // 新增訊息
		messageGroup.POST("/delete", mr.message.Delete) // 刪除訊息
		messageGroup.POST("/update", mr.message.Update) // 更新訊息
		// messageGroup.POST("/query", mr.message.Create)  // 查找訊息
	}
}
