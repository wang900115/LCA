package controller

import (
	response "LCA/internal/adapter/gin/controller/response/json"
	"LCA/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	response response.JSONResponse
	user     usecase.UserUsecase
}

func NewUserController(response response.JSONResponse, user usecase.UserUsecase) *UserController {
	return &UserController{response: response, user: user}
}

func (uc *UserController) CreateUser(c *gin.Context) {

}

func (uc *UserController) DeleteUser(c *gin.Context) {

}
