package route

import (
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/controller"
	middleware "github.com/wang900115/LCA/pkg/common/middleware/jwt"
	"github.com/wang900115/LCA/pkg/common/router"
)

type UserRouter struct {
	user controller.UserController
	jwt  middleware.JWT
}

func NewUserRouter(user *controller.UserController, jwt *middleware.JWT) router.IRoute {
	return &UserRouter{user: *user, jwt: *jwt}
}

func (ur *UserRouter) Setup(router *gin.RouterGroup) {
	userGroup := router.Group("v1/user/")
	{
		userGroup.POST("/login", ur.user.Login)                      // 登入
		userGroup.POST("/register", ur.user.Register)                // 註冊
		userGroup.POST("/logout", ur.jwt.Middleware, ur.user.Logout) // 登出

		userGroup.POST("/query", ur.user.Query)                      // 列出
		userGroup.POST("/update", ur.jwt.Middleware, ur.user.Update) // 更新
		userGroup.POST("/delete", ur.jwt.Middleware, ur.user.Delete) // 刪除
	}
}
