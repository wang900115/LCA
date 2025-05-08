package router

import (
	"LCA/internal/adapter/gin/controller"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController controller.UserController
	jwt            jwt.JWT
}

func NewUserRouter(userController controller.UserController, jwt jwt.JWT) *UserRouter {
	return &UserRouter{userController: userController, jwt: jwt}
}

func (ur *UserRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("v1/User/")
	{
		userGroup.POST("/participate", ur.userController.CreateUser)
		userGroup.DELETE("/leave", ur.jwt.Middleware, ur.userController.DeleteUser)
	}
}
