package controller

import (
	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type ChannelController struct {
	response iresponse.IResponse
	channel  usecase.ChannelUsecase
}

func NewChannelController(reponse iresponse.IResponse, channel *usecase.ChannelUsecase) *ChannelController {
	return &ChannelController{response: reponse, channel: *channel}
}

func (cc *ChannelController) CreateChannel(c *gin.Context) {
	var request validator.ChannelCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		cc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
	}

	err := cc.channel.CreateChannel(c, request)
	if err != nil {
		cc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	cc.response.Success(c, SUCCESS)
}

func (cc *ChannelController) QueryUsers(c *gin.Context) {
	var request validator.ChannelQueryUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		cc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
	}
	users, err := cc.channel.QueryChannelUsers(c, request.ChannelID)
	if err != nil {
		cc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	cc.response.SuccessWithData(c, ACCEPTED_SUCCESS, users)
}

func (cc *ChannelController) QueryMessages(c *gin.Context) {
	var request validator.ChannelQueryUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		cc.response.ValidatorFail(c, INVALID_PARAM_ERROR)
	}
	channels, err := cc.channel.QueryChannelMessages(c, request.ChannelID)
	if err != nil {
		cc.response.FailWithError(c, COMMON_INTERNAL_ERROR, err)
		return
	}
	cc.response.SuccessWithData(c, ACCEPTED_SUCCESS, channels)
}
