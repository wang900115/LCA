package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
)

type ChannelUserController struct {
	channel usecase.ChannelUsecase
	resp    iresponse.IResponse
}

func NewChannelUserController(channel *usecase.ChannelUsecase, resp iresponse.IResponse) *ChannelUserController {
	return &ChannelUserController{channel: *channel, resp: resp}
}

func (cu *ChannelUserController) GetChannelUsers(c *gin.Context) {
	var req validator.GetChannelUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cu.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	users, err := cu.channel.GetChannelUsers(c, req.ChannelID)
	if err != nil {
		cu.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	cu.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"users": users,
	})
}
