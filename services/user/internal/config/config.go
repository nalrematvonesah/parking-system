package config

import "os"

type Config struct {
	GRPCPort  string
	DBURL     string
	RedisAddr string
	RedisDB   int
}

func Load() Config {
	return Config{
		GRPCPort:  get("GRPC_PORT", "50052"),
		DBURL:     get("DATABASE_URL", "postgres://postgres:123@postgres:5432/parking?sslmode=disable"),
		RedisAddr: get("REDIS_ADDR", "redis:6379"),
		RedisDB:   1,
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
