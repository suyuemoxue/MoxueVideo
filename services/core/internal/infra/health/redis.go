package health

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

type RedisChecker struct {
	client *redis.Client
}

func NewRedisChecker(client *redis.Client) *RedisChecker {
	return &RedisChecker{client: client}
}

func (c *RedisChecker) Check(ctx context.Context) error {
	if c == nil || c.client == nil {
		return errors.New("redis unavailable")
	}
	return c.client.Ping(ctx).Err()
}
