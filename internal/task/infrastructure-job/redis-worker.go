package infrastructurejob

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisJob struct {
	logger *zap.Logger
	redis  *redis.Client
}

func NewRedisJob(logger *zap.Logger, redis *redis.Client) *RedisJob {
	return &RedisJob{logger: logger, redis: redis}
}

func (rj *RedisJob) SetUp(s gocron.Scheduler) {
	_, err := s.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(rj.Health))
	if err != nil {
		rj.logger.Error(err.Error(), zap.String("action", "[setup]infrastruction-redis-health"))
	}

	_, err = s.NewJob(gocron.DurationJob(24*time.Hour), gocron.NewTask(rj.Clean))
	if err != nil {
		rj.logger.Error(err.Error(), zap.String("action", "[setup]infrastruction-redis-clean"))
	}
}

func (rj *RedisJob) Clean() error {
	ctx := context.Background()
	iter := rj.redis.Scan(ctx, 0, "user", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := rj.redis.Del(ctx, key).Err(); err != nil {
			rj.logger.Error(err.Error(), zap.String("action", "infrastruction-redis-clean"))
		}
	}
	return iter.Err()
}

func (rj *RedisJob) Health() {
	ctx := context.Background()
	if err := rj.redis.Ping(ctx).Err(); err != nil {
		rj.logger.Error(err.Error(), zap.String("action", "infrastruction-redis-health"))
	}
}
