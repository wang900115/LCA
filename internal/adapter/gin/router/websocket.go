package router

import (
	"LCA/internal/adapter/gin/controller"

	"github.com/gin-gonic/gin"
)

type WebSocketRouter struct {
	webSocketController controller.WebSocketController
}

func NewWebSocketRouter(webSocketController *controller.WebSocketController) IRoute {
	return &WebSocketRouter{webSocketController: *webSocketController}
}

func (wr *WebSocketRouter) Setup(router *gin.RouterGroup) {
	webSocketGroup := router.Group("v1/ws/")
	{
		webSocketGroup.GET("/connect", wr.webSocketController.Handle)
	}
}
