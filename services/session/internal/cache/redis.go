package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func New(addr string) *Redis {
	return &Redis{
		rdb: redis.NewClient(
			&redis.Options{
				Addr: addr,
			},
		),
	}
}

func (r *Redis) SetSession(
	ctx context.Context,
	key string,
	value string,
) error {
	return r.rdb.Set(
		ctx,
		key,
		value,
		time.Hour,
	).Err()
}

func (r *Redis) GetSession(
	ctx context.Context,
	key string,
) (string, error) {
	return r.rdb.Get(
		ctx,
		key,
	).Result()
}
