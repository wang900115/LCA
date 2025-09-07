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
		cc.response.ValidatorFail(c, validatorFail)
	}

	err := cc.channel.CreateChannel(c, request)
	if err != nil {
		cc.response.FailWithError(c, createFail, err)
		return
	}
	cc.response.Success(c, createSuccess)
}

func (cc *ChannelController) QueryUsers(c *gin.Context) {
	var request validator.ChannelQueryUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		cc.response.ValidatorFail(c, validatorFail)
	}
	users, err := cc.channel.QueryChannelUsers(c, request.ChannelID)
	if err != nil {
		cc.response.FailWithError(c, queryFail, err)
		return
	}
	cc.response.SuccessWithData(c, querySuccess, users)
}

func (cc *ChannelController) QueryMessage(c *gin.Context) {
	var request validator.ChannelQueryUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		cc.response.ValidatorFail(c, validatorFail)
	}
	channels, err := cc.channel.QueryChannelMessages(c, request.ChannelID)
	if err != nil {
		cc.response.FailWithError(c, queryFail, err)
		return
	}
	cc.response.SuccessWithData(c, querySuccess, channels)
}
