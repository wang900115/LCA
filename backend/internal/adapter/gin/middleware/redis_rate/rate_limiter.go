package redisrate

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Option struct {
	LimitPerMinute int `yaml:"limit_per_minute"`
}

func NewOption(conf *viper.Viper) Option {
	return Option{
		LimitPerMinute: conf.GetInt("ratelimit.limit_per_minute"),
	}
}

type RateLimiter struct {
	limiter *redis_rate.Limiter
	logger  *zap.Logger

	LimitPerMinute int
}

func NewRateLimiter(redisPool *redis.Client, logger *zap.Logger, option Option) *RateLimiter {
	return &RateLimiter{
		limiter:        redis_rate.NewLimiter(redisPool),
		logger:         logger,
		LimitPerMinute: option.LimitPerMinute,
	}
}

func (rl RateLimiter) Middleware(c *gin.Context) {
	clientIP := c.ClientIP()
	if clientIP != "" {
		limitPerMinute := 600
		if rl.LimitPerMinute > 0 {
			limitPerMinute = rl.LimitPerMinute
		}

		res, err := rl.limiter.Allow(c, clientIP, redis_rate.PerSecond(limitPerMinute))
		if err != nil {
			rl.logger.Error(err.Error(),
				zap.String("type", "rate limiter error"))
			c.Abort()
			return
		}

		h := c.Writer.Header()
		h.Set("X-RateLimit-Limit", strconv.FormatInt(int64(rl.LimitPerMinute), 10))
		h.Set("X-RateLimit-Remaining", strconv.FormatInt(int64(res.Remaining), 10))
		h.Set("X-RateLimit-Delay", strconv.FormatInt(int64(res.ResetAfter/time.Second), 10))

		if res.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": fmt.Sprintf("retry after %d seconds", res.RetryAfter/time.Second),
			})

			c.Abort()
			return
		}
	}

	c.Next()
}
