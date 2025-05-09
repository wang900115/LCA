package controller

import (
	iresponse "LCA/internal/adapter/gin/controller/response"
	"LCA/internal/adapter/gin/validator"
	"LCA/internal/application/usecase"

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
	userUUID, err := uc.user.CreateUser(request.ChannelUUID, request.Username)
	if err != nil {
		uc.response.FailWithError(c, createFail, err)
		return
	}
	token, err := uc.token.CreateToken(userUUID, request.ChannelUUID, request.Username)
	if err != nil {
		uc.response.FailWithError(c, createFail, err)
		return
	}
	uc.response.SuccessWithData(c, createSuccess, token)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	UserUUID := c.GetString("user_uuid")
	userUUID, err := uc.user.DeleteUser(UserUUID)
	if err != nil {
		uc.response.FailWithError(c, deleteFail, err)
		return
	}

	uc.response.SuccessWithData(c, deleteSuccess, userUUID)
}
