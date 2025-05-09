package router

import (
	"LCA/internal/adapter/gin/controller"
	"LCA/internal/adapter/gin/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type WebSocketRouter struct {
	webSocketController controller.WebSocketController
	jwt                 jwt.JWT
}

func NewWebSocketRouter(webSocketController *controller.WebSocketController, jwt *jwt.JWT) IRoute {
	return &WebSocketRouter{webSocketController: *webSocketController, jwt: *jwt}
}

func (wr *WebSocketRouter) Setup(router *gin.RouterGroup) {
	webSocketGroup := router.Group("v1/ws/")
	{
		webSocketGroup.GET("/connect", wr.jwt.Middleware, wr.webSocketController.Handle)
	}
}
