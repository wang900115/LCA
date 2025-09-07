package controller

import (
	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	response iresponse.IResponse
	message  usecase.MessageUsecase
}

func NewMessageController(response iresponse.IResponse, message *usecase.MessageUsecase) *MessageController {
	return &MessageController{response: response, message: *message}
}

func (mc *MessageController) ReadMessage(c *gin.Context) {
	var request validator.MessageReadRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
		return
	}
	message, err := mc.message.ReadMessage(c, request.MessageId)
	if err != nil {
		mc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	mc.response.SuccessWithData(c, ACCEPTED_SUCCESS, message)
}
