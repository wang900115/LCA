package router

import (
	"LCA/internal/adapter/gin/controller"

	"github.com/gin-gonic/gin"
)

type MessageRouter struct {
	messageController controller.MessageController
	jwt               jwt.JWT
}

func NewMessageRouter(messageController controller.MessageController, jwt jwt.JWT) *MessageRouter {
	return &MessageRouter{messageController: messageController, jwt: jwt}
}

func (mr *MessageRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("v1/message/")
	{
		userGroup.POST("/create", mr.jwt.Middleware, mr.messageController.CreateMessage)
		userGroup.DELETE("/delete", mr.jwt.Middleware, mr.messageController.DeleteMessage)
		userGroup.GET("/query", mr.jwt.Middleware, mr.messageController.QueryMessage)
	}
}
