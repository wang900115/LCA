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

func (mc *MessageController) CreateMessage(c *gin.Context) {

	ChannelUUID := c.GetString("channel_uuid")
	UserUUID := c.GetString("user_uuid")
	var request validator.MessageCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		mc.response.ValidatorFail(c, validatorFail)
		return
	}
	messageUUID, err := mc.message.CreateMessage(ChannelUUID, UserUUID, request.Content)
	if err != nil {
		mc.response.FailWithError(c, createFail, err)
		return
	}
	mc.response.SuccessWithData(c, createSuccess, messageUUID)
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
}

func (mc *MessageController) QueryMessage(c *gin.Context) {
	channelUUID := c.GetString("channel_uuid")
	messages, err := mc.message.QueryMessages(channelUUID)
	if err != nil {
		mc.response.FailWithError(c, queryFail, err)
		return
	}
	mc.response.SuccessWithData(c, querySuccess, messages)
}
