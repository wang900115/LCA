package cors

import (
	"net/http"

	"github.com/wang900115/LCA/internal/adapter/gin/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Option struct {
	AllowOrigin string `yaml:"allow_origin"`
}

func NewOption(conf *viper.Viper) Option {
	return Option{
		AllowOrigin: conf.GetString("server.allow_origin"),
	}
}

type CORS struct {
	allowOrigin string
}

func NewCORS(option Option) middleware.IMiddleware {
	return &CORS{allowOrigin: option.AllowOrigin}
}

func (cors CORS) Middleware(c *gin.Context) {
	method := c.Request.Method
	origin := c.Request.Header.Get("Origin")

	if origin != "" {
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Origin", cors.allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
	}

	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
