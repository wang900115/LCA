package controller

import (
	"net/http"

	iresponse "github.com/wang900115/LCA/internal/adapter/gin/controller/response"
	"github.com/wang900115/LCA/internal/adapter/websocket/connection"
	"github.com/wang900115/LCA/internal/application/usecase"

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

	userName, _ := c.GetQuery("username")
	channelName, _ := c.GetQuery("channelname")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsc.response.FailWithError(c, accessDenied, err)
		return
	}

	client := connection.NewClient(userName, channelName, conn, &wsc.message)
	wsc.hub.Register <- client

	go client.ReadPump(&wsc.hub)
	go client.WritePump()
}
