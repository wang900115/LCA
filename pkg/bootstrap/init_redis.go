package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ReadRedis struct {
	Redis *redis.Client
	Label string
}

type RedisGroup struct {
	Write *redis.Client
	Reads []ReadRedis
}

func (c RedisConfig) DSN() *redis.Options {
	return &redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	}
}

func NewRedisGroup(v *viper.Viper) *RedisGroup {
	cfg := RedisConfig{
		Addr:     v.GetString("redis.write.addr"),
		Password: v.GetString("redis.write.password"),
		DB:       v.GetInt("redis.write.db"),
	}

	write := redis.NewClient(cfg.DSN())

	var reads []ReadRedis
	readConfigs := v.Get("redis.reads").([]interface{})

	for _, cfg := range readConfigs {
		rc := cfg.(map[string]interface{})

		conf := RedisConfig{
			Addr:     rc["addr"].(string),
			Password: rc["password"].(string),
			DB:       rc["db"].(int),
		}

		read := redis.NewClient(conf.DSN())
		reads = append(reads, ReadRedis{Redis: read, Label: conf.Addr})
	}

	return &RedisGroup{Write: write, Reads: reads}
}

func (rg *RedisGroup) PickLeastConnRead() *redis.Client {
	var min *redis.Client
	var minConns int

	for _, r := range rg.Reads {
		stats := r.Redis.PoolStats()
		if min == nil || stats.TotalConns < uint32(minConns) {
			min = r.Redis
			minConns = int(stats.TotalConns)
		}
	}

	return min
}
