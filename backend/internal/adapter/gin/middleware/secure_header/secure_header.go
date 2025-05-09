package secureheader

import (
	"github.com/gin-gonic/gin"
)

type SecureHeader struct {
}

func NewSecureHeader() *SecureHeader {
	return &SecureHeader{}
}

func (SecureHeader) Middleware(c *gin.Context) {
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Xss-Protection", "1; mod=block")
	c.Header("Referrer-Policy", "no-referrer")

	c.Next()
}
