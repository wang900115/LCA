package middleware

import "github.com/gin-gonic/gin"

type IMiddleware interface {
	Middleware(*gin.Context)
}
