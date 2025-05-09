package controller

import (
	iresponse "LCA/internal/adapter/gin/controller/response"
	"LCA/internal/adapter/gin/validator"
	"LCA/internal/application/usecase"

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
	channelUUID, err := cc.channel.CreateChannel()
	if err != nil {
		cc.response.FailWithError(c, createFail, err)
		return
	}
	cc.response.SuccessWithData(c, createSuccess, channelUUID)
}

func (cc *ChannelController) QueryUsers(c *gin.Context) {
	var request validator.ChannelQueryUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		cc.response.ValidatorFail(c, validatorFail)
		return
	}

	users, err := cc.channel.QueryUsers(request.ChannelUUID)
	if err != nil {
		cc.response.FailWithError(c, queryFail, err)
		return
	}
	cc.response.SuccessWithData(c, querySuccess, users)
}

func (cc *ChannelController) QueryChannel(c *gin.Context) {
	channels, err := cc.channel.QueryChannels()
	if err != nil {
		cc.response.FailWithError(c, queryFail, err)
		return
	}
	cc.response.SuccessWithData(c, querySuccess, channels)
}
