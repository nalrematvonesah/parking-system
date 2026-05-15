package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort    string
	DBURL       string
	RedisAddr   string
	NatsURL     string
	ParkingAddr string
	// PricePerHour is the base price for one billable hour.
	PricePerHour float64
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	return Config{
		GRPCPort:     get("GRPC_PORT", "50053"),
		DBURL:        get("DATABASE_URL", "postgres://postgres:123@postgres:5432/parking?sslmode=disable"),
		RedisAddr:    get("REDIS_ADDR", "redis:6379"),
		NatsURL:      get("NATS_URL", "nats://nats:4222"),
		ParkingAddr:  get("PARKING_ADDR", "parking:50051"),
		PricePerHour: 500.0,
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
