package router

import (
	"github.com/wang900115/LCA/internal/adapter/controller"
	"github.com/wang900115/LCA/internal/adapter/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController controller.UserController
	jwt            jwt.JWT
}

func NewUserRouter(userController *controller.UserController, jwt *jwt.JWT) IRoute {
	return &UserRouter{userController: *userController, jwt: *jwt}
}

func (ur *UserRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("v1/User/")
	{
		userGroup.POST("/participate", ur.userController.CreateUser)
		userGroup.DELETE("/leave", ur.jwt.Middleware, ur.userController.DeleteUser)
	}
}
