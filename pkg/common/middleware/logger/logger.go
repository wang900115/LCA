package logger

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l Logger) Middleware(c *gin.Context) {
	var requestBody string
	if c.Request != nil {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			l.logger.Error(err.Error(),
				zap.String("type", "Logger Middleware error"))
		}

		requestBody = string(body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	c.Next()

	l.logger.Info("Middleware logger",
		zap.String("IP_From", c.ClientIP()),
		zap.String("method", c.Request.Method),
		zap.Int("status", c.Writer.Status()),
		zap.String("url", c.Request.RequestURI),
		zap.String("request", requestBody),
	)

	c.Next()
}
