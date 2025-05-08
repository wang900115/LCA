package controller

import (
	response "LCA/internal/adapter/gin/controller/response/json"
	"LCA/internal/adapter/gin/validator"
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
	var request validator.MessageCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mc.response.ValidatorFail(c, validatorFail)
		return
	}
	messageUUID, err := mc.message.CreateMessage(request.UserUUID, request.ChannelUUID, request.Content)
	if err != nil {
		mc.response.FailWithError(c, createFail, err)
		return
	}
	mc.response.SuccessWithData(c, createSuccess, messageUUID)
	return
}

func (mc *MessageController) DeleteMessage(c *gin.Context) {
	var request validator.MessageDeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mc.response.ValidatorFail(c, validatorFail)
		return
	}
	if err := mc.message.DeleteMessage(request.MessageUUID); err != nil {
		mc.response.FailWithError(c, deleteFail, err)
		return
	}
	mc.response.Success(c, deleteSuccess)
	return
}

func (mc *MessageController) QueryMessage(c *gin.Context) {
	var request validator.MessageQueryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mc.response.ValidatorFail(c, validatorFail)
		return
	}
	messages, err := mc.message.QueryMessages(request.ChannelUUID)
	if err != nil {
		mc.response.FailWithError(c, queryFail, err)
		return
	}
	mc.response.SuccessWithData(c, querySuccess, messages)
	return
}
