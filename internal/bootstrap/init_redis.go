package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type redisoption struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewRedisOption(conf *viper.Viper) redisoption {
	return redisoption{
		Addr:     conf.GetString("redis.host"),
		Username: conf.GetString("redis.user"),
		Password: conf.GetString("redis.password"),
		DB:       conf.GetInt("redis.database"),
	}
}

func NewRedisPool(option redisoption) *redis.Client {
	redisPool := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Username: option.Username,
		Password: option.Password,
		DB:       option.DB,
	})
	return redisPool
}
