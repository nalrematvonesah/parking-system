package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func New(addr string, db int) *RedisCache {
	return &RedisCache{
		rdb: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   db,
		}),
	}
}

const keyAvailable = "parking:available"

func (c *RedisCache) GetAvailable(ctx context.Context) (int32, error) {
	val, err := c.rdb.Get(ctx, keyAvailable).Int()
	return int32(val), err
}

func (c *RedisCache) SetAvailable(ctx context.Context, v int32) {
	c.rdb.Set(ctx, keyAvailable, v, time.Minute)
}

func (c *RedisCache) InvalidateAvailable(ctx context.Context) {
	c.rdb.Del(ctx, keyAvailable)
}
