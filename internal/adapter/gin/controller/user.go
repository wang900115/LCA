package controller

import (
	iresponse "github.com/wang900115/LCA/internal/adapter/gin/controller/response"
	"github.com/wang900115/LCA/internal/adapter/gin/validator"
	"github.com/wang900115/LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	response iresponse.IResponse
	token    usecase.TokenUsecase
	user     usecase.UserUsecase
}

func NewUserController(response iresponse.IResponse, token *usecase.TokenUsecase, user *usecase.UserUsecase) *UserController {
	return &UserController{response: response, token: *token, user: *user}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var request validator.UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, validatorFail)
		return
	}
	userName, err := uc.user.CreateUser(request.ChannelName, request.Username)
	if err != nil {
		uc.response.FailWithError(c, createFail, err)
		return
	}

	token, err := uc.token.CreateToken(userName, request.ChannelName)
	if err != nil {
		uc.response.FailWithError(c, createFail, err)
		return
	}
	uc.response.SuccessWithData(c, createSuccess, token)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	User := c.GetString("user")
	Channel := c.GetString("channel")

	userName, err := uc.user.DeleteUser(User)
	if err != nil {
		uc.response.FailWithError(c, deleteFail, err)
		return
	}

	uc.token.DeleteToken(User, Channel)
	uc.response.SuccessWithData(c, deleteSuccess, userName)
}
