package controller

import (
	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
	"github.com/wang900115/LCA/internal/adapter/validator"
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

func (uc *UserController) Register(c *gin.Context) {
	var request validator.UserCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, validatorFail)
		return
	}
	err := uc.user.CreateUser(c, request)
	if err != nil {
		uc.response.FailWithError(c, createFail, err)
		return
	}
	uc.response.Success(c, createSuccess)
}

func (uc *UserController) Delete(c *gin.Context) {
	id := c.GetUint("user_id")
	if err := uc.user.DeleteUser(c, id); err != nil {
		uc.response.FailWithError(c, deleteFail, err)
		return
	}
	uc.response.Success(c, deleteSuccess)
	return
}

func (uc *UserController) Login(c *gin.Context) {
	var request validator.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		uc.response.ValidatorFail(c, validatorFail)
		return
	}
	id, userLogin, err := uc.user.Login(c, request)
	if id == nil || userLogin == nil || err != nil {
		uc.response.FailWithError(c, accessDenied, err)
		return
	}
	token, err := uc.token.UserLoginGenerateToken(*id, *userLogin)
	if err != nil {
		uc.response.FailWithError(c, accessDenied, err)
		return
	}
	uc.response.SuccessWithData(c, querySuccess, token)
}

func (uc *UserController) Logout(c *gin.Context) {
	ip := c.GetString("ip_address")
	if ip != c.ClientIP() {
		uc.response.Fail(c, accessDenied)
		return
	}
	id, err := uc.user.Logout(c, c.ClientIP())
	if id == nil || err != nil {
		uc.response.FailWithError(c, accessDenied, err)
		return
	}
	if err := uc.token.DeleteUserToken(*id); err != nil {
		uc.response.FailWithError(c, accessDenied, err)
		return
	}
	uc.response.Success(c, querySuccess)
}
