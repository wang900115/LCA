package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
)

type ChannelMessageController struct {
	message usecase.MessageUsecase
	resp    iresponse.IResponse
}

func NewChannelMessageController(message *usecase.MessageUsecase, resp iresponse.IResponse) *ChannelMessageController {
	return &ChannelMessageController{message: *message, resp: resp}
}

func (cm *ChannelMessageController) GetChannelMessages(c *gin.Context) {
	var req validator.GetChannelMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		cm.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	messages, err := cm.message.GetChannelMessages(c, req.ChannelID)
	if err != nil {
		cm.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	cm.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"messages": messages,
	})
}
