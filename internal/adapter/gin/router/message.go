package router

import (
	"github.com/wang900115/LCA/internal/adapter/gin/controller"
	"github.com/wang900115/LCA/internal/adapter/gin/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type MessageRouter struct {
	messageController controller.MessageController
	jwt               jwt.JWT
}

func NewMessageRouter(messageController *controller.MessageController, jwt *jwt.JWT) IRoute {
	return &MessageRouter{messageController: *messageController, jwt: *jwt}
}

func (mr *MessageRouter) Setup(router *gin.RouterGroup) {
	messageGroup := router.Group("v1/message/")
	{
		messageGroup.POST("/create", mr.jwt.Middleware, mr.messageController.CreateMessage)
		messageGroup.DELETE("/delete", mr.jwt.Middleware, mr.messageController.DeleteMessage)
		messageGroup.POST("/query", mr.jwt.Middleware, mr.messageController.QueryMessage)
	}

}
