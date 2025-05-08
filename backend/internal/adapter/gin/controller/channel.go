package controller

import (
	response "LCA/internal/adapter/gin/controller/response/json"
	"LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type ChannelController struct {
	response response.JSONResponse
	channel  usecase.ChannelUsecase
}

func NewChannelController(reponse response.JSONResponse, channel usecase.ChannelUsecase) *ChannelController {
	return &ChannelController{response: reponse, channel: channel}
}

func (cc *ChannelController) CreateChannel(c *gin.Context) {

}

func (cc *ChannelController) DeleteChannel(c *gin.Context) {

}

func (cc *ChannelController) QueryChannel(c *gin.Context) {

}
