package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/validator"
	"github.com/wang900115/LCA/internal/application/usecase"
	"github.com/wang900115/LCA/pkg/common"
	iresponse "github.com/wang900115/LCA/pkg/common/response"
	"github.com/wang900115/LCA/pkg/domain"
)

type MessageController struct {
	message usecase.MessageUsecase
	resp    iresponse.IResponse
}

func NewMessageController(message *usecase.MessageUsecase, resp iresponse.IResponse) *MessageController {
	return &MessageController{message: *message, resp: resp}
}

func (m *MessageController) Create(c *gin.Context) {
	var req validator.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	userID := c.GetUint("user_id")
	message := domain.Message{
		ChannelID: req.ChannelID,
		UserID:    userID,
		MsgType:   req.MsgType,
		Status:    req.Status,
		ReplyToID: req.ReplyToID,

		Content:   req.Content,
		AttachURL: req.AttachURL,
	}
	created, err := m.message.CreateMessage(c, message)
	if err != nil {
		m.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	m.resp.SuccessWithData(c, common.CREATE_SUCCESS, map[string]interface{}{
		"message": created,
	})
}

func (m *MessageController) Update(c *gin.Context) {
	var req validator.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}
	userID := c.GetUint("user_id")
	message := domain.Message{
		ChannelID: req.ChannelID,
		UserID:    userID,
		MsgType:   req.MsgType,
		Status:    req.Status,
		ReplyToID: req.ReplyToID,

		Content:   req.Content,
		AttachURL: req.AttachURL,
	}
	updated, err := m.message.UpdateMessage(c, message)
	if err != nil {
		m.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	m.resp.SuccessWithData(c, common.UPDATE_SUCCESS, map[string]interface{}{
		"message": updated,
	})
}

func (m *MessageController) Delete(c *gin.Context) {
	var req validator.DeleteMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		m.resp.FailWithError(c, common.PARAM_ERROR, err)
		return
	}

	err := m.message.DeleteMessage(c, req.MeesageID)
	if err != nil {
		m.resp.FailWithError(c, common.INTERNAL_SERVICE_ERROR, err)
		return
	}

	m.resp.Success(c, common.DELETE_SUCCESS)
}
