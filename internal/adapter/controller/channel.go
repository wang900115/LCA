package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
	"github.com/wang900115/LCA/pkg/domain"
)

type ChannelController struct {
	channel usecase.ChannelUsecase
	resp    iresponse.IResponse
}

func NewChannelController(channel *usecase.ChannelUsecase, resp iresponse.IResponse) *ChannelController {
	return &ChannelController{channel: *channel, resp: resp}
}

func (cc *ChannelController) GetAllChannels(c *gin.Context) {
	channels, err := cc.channel.GetAllChannels(c)
	if err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cc.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"channels": channels,
	})
}

func (cc *ChannelController) Create(c *gin.Context) {
	var req validator.CreateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	channel := domain.Channel{
		ChannelName: req.ChannelName,
		ChannelType: req.ChannelType,
	}

	created, err := cc.channel.CreateChannel(c, channel)
	if err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	cc.resp.SuccessWithData(c, common.CREATE_SUCCESS, map[string]interface{}{
		"channel": created,
	})
}
func (cc *ChannelController) Update(c *gin.Context) {
	var req validator.UpdateChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	channel := domain.Channel{
		ID:          req.ID,
		ChannelName: req.ChannelName,
		ChannelType: req.ChannelType,
	}

	updated, err := cc.channel.UpdateChannel(c, channel)
	if err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	cc.resp.SuccessWithData(c, common.UPDATE_SUCCESS, map[string]interface{}{
		"channel": updated,
	})
}

func (cc *ChannelController) Delete(c *gin.Context) {
	var req validator.DeleteChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cc.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	if err := cc.channel.DeleteChannel(c, req.ChannelID); err != nil {
		cc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cc.resp.Success(c, common.DELETE_SUCCESS)
}
