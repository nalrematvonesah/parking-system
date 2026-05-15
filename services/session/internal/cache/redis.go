package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func New(addr string) *Redis {
	return &Redis{
		rdb: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

// SetActiveSession caches that a given session is active. TTL keeps it cheap
// even if a session is never closed.
func (r *Redis) SetActiveSession(ctx context.Context, sessionID int64) error {
	return r.rdb.Set(ctx, keyActive(sessionID), "1", 24*time.Hour).Err()
}

func (r *Redis) IsActive(ctx context.Context, sessionID int64) (bool, error) {
	n, err := r.rdb.Exists(ctx, keyActive(sessionID)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *Redis) DeleteActive(ctx context.Context, sessionID int64) error {
	return r.rdb.Del(ctx, keyActive(sessionID)).Err()
}

func keyActive(id int64) string {
	return fmt.Sprintf("session:active:%s", strconv.FormatInt(id, 10))
}
