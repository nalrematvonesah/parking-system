package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func New(addr string, db int) *RedisCache {
	return &RedisCache{
		rdb: redis.NewClient(&redis.Options{Addr: addr, DB: db}),
	}
}

func userKey(id int64) string { return fmt.Sprintf("user:session:%d", id) }

func (c *RedisCache) SetSession(ctx context.Context, userID int64) {
	c.rdb.Set(ctx, userKey(userID), "active", 24*time.Hour)
}

func (c *RedisCache) InvalidateSession(ctx context.Context, userID int64) {
	c.rdb.Del(ctx, userKey(userID))
}
