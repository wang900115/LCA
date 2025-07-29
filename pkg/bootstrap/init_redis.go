package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/wang900115/LCA/pkg/common"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type GenernalRedis struct {
	Redis *redis.Client
	Label string
}

type SentinelRedis struct {
	Redis *redis.SentinelClient
	Label string
}

type RedisGroup struct {
	Write      *GenernalRedis
	Reads      []*GenernalRedis
	Sentinels  []*SentinelRedis
	MasterName string
}

func (c RedisConfig) DSN() *redis.Options {
	return &redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	}
}

func loadSentinels(v *viper.Viper) []*SentinelRedis {
	var sentinels []*SentinelRedis
	sentinelConfigs := v.Get("redis.sentinel").([]interface{})

	for _, cfg := range sentinelConfigs {
		rc := cfg.(map[string]interface{})

		conf := RedisConfig{
			Addr:     rc["addr"].(string),
			Password: rc["password"].(string),
			DB:       rc["db"].(int),
		}

		client := redis.NewSentinelClient(conf.DSN())
		sentinels = append(sentinels, &SentinelRedis{Redis: client, Label: conf.Addr})
	}

	return sentinels
}

func getMasterAddrFromSentinel(ctx context.Context, sentinels []*SentinelRedis, masterName string) (string, error) {
	for _, s := range sentinels {
		addr, err := s.Redis.GetMasterAddrByName(ctx, masterName).Result()
		if err != nil {
			continue
		}
		if len(addr) == 2 {
			return fmt.Sprintf("%s:%s", addr[0], addr[1]), nil
		}
	}
	return "", common.RedisSentinelMaster
}

func NewRedisGroup(v *viper.Viper) *RedisGroup {
	ctx := context.Background()
	sentinels := loadSentinels(v)
	masterName := v.GetString("redis.master_name")

	var writeClient *redis.Client
	var writeLabel string

	if len(sentinels) > 0 && masterName != "" {
		addr, err := getMasterAddrFromSentinel(ctx, sentinels, masterName)
		if err != nil {
			log.Fatalf("failed to get master sentinel: %v", err)
		}
		opt := &redis.Options{
			Addr:     addr,
			Password: v.GetString("redis.write.password"),
			DB:       v.GetInt("redis.write.db"),
		}
		writeClient = redis.NewClient(opt)
		writeLabel = addr
	} else {
		opt := &RedisConfig{
			Addr:     v.GetString("redis.write.addr"),
			Password: v.GetString("redis.write.password"),
			DB:       v.GetInt("redis.write.db"),
		}
		writeClient = redis.NewClient(opt.DSN())
		writeLabel = opt.Addr
	}

	var reads []*GenernalRedis
	readConfigs := v.Get("redis.reads").([]interface{})
	for _, cfg := range readConfigs {
		rc := cfg.(map[string]interface{})
		conf := RedisConfig{
			Addr:     rc["addr"].(string),
			Password: rc["password"].(string),
			DB:       rc["db"].(int),
		}
		client := redis.NewClient(conf.DSN())
		reads = append(reads, &GenernalRedis{Redis: client, Label: conf.Addr})
	}

	return &RedisGroup{
		Write:      &GenernalRedis{Redis: writeClient, Label: writeLabel},
		Reads:      reads,
		Sentinels:  sentinels,
		MasterName: masterName,
	}
}

func (re *RedisGroup) PickRedisLeastConnRead() *redis.Client {
	var min *redis.Client
	var minConns int

	for _, r := range re.Reads {
		stats := r.Redis.PoolStats()
		if min == nil || stats.TotalConns < uint32(minConns) {
			min = r.Redis
			minConns = int(stats.TotalConns)
		}
	}

	return min
}

func (re *RedisGroup) HeadlthCheck(ctx context.Context) {
	if err := pingRedisNode(ctx, re.Write.Redis); err != nil {
		log.Printf("[Failover-Redis-Warning] write redis failed: %v", err)

		addr, err := getMasterAddrFromSentinel(ctx, re.Sentinels, re.MasterName)
		if err != nil {
			log.Printf("[Failover-Redis-Error] sentinel lookup failed: %v", err)
			return
		}
		newWrite := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: re.Write.Redis.Options().Password,
			DB:       re.Write.Redis.Options().DB,
		})

		if pingRedisNode(ctx, newWrite) == nil {
			log.Printf("[Failover-Redis-Info] switched to new master: %s", addr)
			re.Write = &GenernalRedis{Redis: newWrite, Label: addr}
		} else {
			log.Printf("[Failover-Redis-Error] new master not healthy: %s", addr)
		}
	}
}

func pingRedisNode(ctx context.Context, redis *redis.Client) error {
	if err := redis.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}
