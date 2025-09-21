package router

import (
	"github.com/wang900115/LCA/internal/adapter/controller"
	"github.com/wang900115/LCA/internal/adapter/middleware/casbin"
	"github.com/wang900115/LCA/internal/adapter/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userController controller.UserController
	userJWT        jwt.USERJWT
	channelJWT     jwt.CHANNELJWT
	casbin         casbin.CASBIN
}

func NewUserRouter(userController *controller.UserController, userJWT *jwt.USERJWT, CHANNELJWT *jwt.CHANNELJWT, casbin *casbin.CASBIN) IRoute {
	return &UserRouter{userController: *userController, userJWT: *userJWT, channelJWT: *CHANNELJWT, casbin: *casbin}
}

func (ur *UserRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("v1/user/")
	{
		userGroup.POST("register", ur.userController.Register)
		userGroup.POST("delete", ur.userJWT.Middleware, ur.userController.Delete)
	}
	userAuthGroup := router.Group("v1/user/auth")
	{
		userAuthGroup.POST("login", ur.userController.Login)
		userAuthGroup.POST("logout", ur.userJWT.Middleware, ur.userController.Logout)
	}
	userJoinGroup := router.Group("v1/user/channel", ur.userJWT.Middleware)
	{
		userJoinGroup.POST("first/join", ur.userController.FirstJoin)
		userJoinGroup.POST("join", ur.userController.Join)
		userJoinGroup.POST("leave", ur.channelJWT.Middleware, ur.userController.LeaveChannel)
	}
	userSpeakGroup := router.Group("v1/user/channel/message", ur.userJWT.Middleware, ur.channelJWT.Middleware)
	{
		userSpeakGroup.POST("comment", ur.userController.Comment)
		userSpeakGroup.POST("edited", ur.userController.Edite)
		userSpeakGroup.POST("delete", ur.userController.Regain)
	}
	userEventGroup := router.Group("v1/user/event", ur.userJWT.Middleware)
	{
		userEventGroup.POST("first/particate")
		userEventGroup.POST("particate")
		userEventGroup.POST("leave")
	}
}
