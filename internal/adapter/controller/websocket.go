package controller

import (
	"net/http"

	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	websocketcore "github.com/wang900115/LCA/internal/adapter/websocket/core"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	hub      websocketcore.Hub
	response iresponse.IResponse
}

func NewWebSocketController(response iresponse.IResponse, hub *websocketcore.Hub) *WebSocketController {
	return &WebSocketController{response: response, hub: *hub}
}

// TODO 要判斷
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsc *WebSocketController) Handle(c *gin.Context) {
	channelId := c.GetUint("channel_id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}

	client := websocketcore.NewClient(conn, channelId)
	wsc.hub.Register <- client

	go client.ReadPump(&wsc.hub)
	go client.WritePump()
}
