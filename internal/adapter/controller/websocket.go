package controller

import (
	"net/http"

	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	websocketcore "github.com/wang900115/LCA/internal/adapter/websocket/core"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	hub      websocketcore.Hub
	user     usecase.UserUsecase
	response iresponse.IResponse
}

func NewWebSocketController(response iresponse.IResponse, user *usecase.UserUsecase, hub *websocketcore.Hub) *WebSocketController {
	return &WebSocketController{response: response, user: *user, hub: *hub}
}

// TODO 要判斷
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsc *WebSocketController) Handle(c *gin.Context) {
	id := c.GetUint("user_id")
	channelId := c.GetUint("channel_id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}

	user, err := wsc.user.ReadUser(c, id)
	if err != nil {
		wsc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}

	client := websocketcore.NewClient(conn, channelId, *user)
	wsc.hub.Register <- client

	go client.ReadPump(&wsc.hub)
	go client.WritePump()
}
