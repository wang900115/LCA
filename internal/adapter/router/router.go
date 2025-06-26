package router

import "github.com/gin-gonic/gin"

type IRoute interface {
	Setup(router *gin.RouterGroup)
}
