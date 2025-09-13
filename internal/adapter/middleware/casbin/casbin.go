package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/wang900115/LCA/internal/adapter/controller"
	iresponse "github.com/wang900115/LCA/internal/adapter/controller/response"
)

type CASBIN struct {
	response iresponse.IResponse
	enforcer *casbin.SyncedEnforcer
}

func NewCASBIN(response iresponse.IResponse, enforcer *casbin.SyncedEnforcer) *CASBIN {
	return &CASBIN{
		response: response,
		enforcer: enforcer,
	}
}

func (cb *CASBIN) Middleware(c *gin.Context) {
	user := c.GetUint("user_id")

	path := c.Request.URL.Path
	method := c.Request.Method
	sub := user

	success, err := cb.enforcer.Enforce(sub, path, method)
	if err != nil {
		cb.response.AuthFail(c, err.Error())
		c.Abort()
		return
	}

	if !success {
		cb.response.AuthFail(c, controller.FORBIDDEN_ERROR)
		c.Abort()
		return
	}

	c.Next()
}
