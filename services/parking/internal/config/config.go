package config

import (
	"os"
)

type Config struct {
	GRPCPort  string
	DBURL     string
	RedisAddr string
	RedisDB   int
}

func Load() Config {
	return Config{
		GRPCPort:  get("GRPC_PORT", "50051"),
		DBURL:     get("DATABASE_URL", "postgres://postgres:123@postgres:5432/parking?sslmode=disable"),
		RedisAddr: get("REDIS_ADDR", "redis:6379"),
		RedisDB:   0,
	}
}

func get(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
