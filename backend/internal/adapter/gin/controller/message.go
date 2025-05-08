package controller

import (
	response "LCA/internal/adapter/gin/controller/response/json"
	"LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	response response.JSONResponse
	message  usecase.MessageUsecase
}

func NewMessageController(response response.JSONResponse, message usecase.MessageUsecase) *MessageController {
	return &MessageController{response: response, message: message}
}

func (mc *MessageController) CreateMessage(c *gin.Context) {

}

func (mc *MessageController) DeleteMessage(c *gin.Context) {

}

func (mc *MessageController) QueryMessage(c *gin.Context) {

}
