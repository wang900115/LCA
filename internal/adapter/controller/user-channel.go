package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
)

type UserChannelController struct {
	channel usecase.ChannelUsecase
	resp    iresponse.IResponse
}

func NewUserChannelController(channel *usecase.ChannelUsecase, resp iresponse.IResponse) *UserChannelController {
	return &UserChannelController{channel: *channel, resp: resp}
}

func (uc *UserChannelController) GetUserChannels(c *gin.Context) {
	userID := c.GetUint("user_id")
	channels, err := uc.channel.GetUserChannels(c, userID)
	if err != nil {
		uc.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}
	uc.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"channels": channels,
	})
}
