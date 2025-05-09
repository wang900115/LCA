package controller

import (
	iresponse "LCA/internal/adapter/gin/controller/response"
	"LCA/internal/adapter/websocket/connection"
	"LCA/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	hub      connection.Hub
	response iresponse.IResponse
	token    usecase.TokenUsecase
	message  usecase.MessageUsecase
}

func NewWebSocketController(response iresponse.IResponse, hub *connection.Hub, token *usecase.TokenUsecase, message *usecase.MessageUsecase) *WebSocketController {
	return &WebSocketController{response: response, hub: *hub, token: *token, message: *message}
}

// TODO 要判斷
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (wsc *WebSocketController) Handle(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	channelUUID := c.GetString("channel_uuid")
	username := c.GetString("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsc.response.FailWithError(c, accessDenied, err)
		return
	}

	client := connection.NewClient(userUUID, channelUUID, username, conn, &wsc.message)
	wsc.hub.Register <- client

	go client.ReadPump(&wsc.hub)
	go client.WritePump()
}
