package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
)

type UserChannelMessageController struct {
	message usecase.MessageUsecase
	resp    iresponse.IResponse
}

func NewUserChannelMessageController(message *usecase.MessageUsecase, resp iresponse.IResponse) *UserChannelMessageController {
	return &UserChannelMessageController{message: *message, resp: resp}
}

func (um *UserChannelMessageController) GetChannelUserMessages(c *gin.Context) {
	var req validator.GetChannelUserMessagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		um.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	messages, err := um.message.GetChannelUserMessages(c, req.ChannelID, req.UserID)
	if err != nil {
		um.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	um.resp.SuccessWithData(c, common.QUERY_SUCCESS, map[string]interface{}{
		"message": messages,
	})
}
