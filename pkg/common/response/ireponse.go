package iresponse

import (
	"github.com/gin-gonic/gin"
)

type IResponse interface {
	// 200
	Success(c *gin.Context, message string)
	// 201
	SuccessWithData(c *gin.Context, message string, data any)

	// 500
	Fail(c *gin.Context, message string)

	// 500
	FailWithError(c *gin.Context, message string, err error)

	// 400
	ValidatorFail(c *gin.Context, message string)

	// 401
	AuthFail(c *gin.Context, message string)
}
