package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort string
	DBURL    string
	NatsURL  string
	SMTP     SMTPConfig
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	return Config{
		GRPCPort: mustGet("GRPC_PORT", "50052"),
		DBURL:    mustGet("DATABASE_URL", "postgres://postgres:123@postgres:5432/parking?sslmode=disable"),
		NatsURL:  mustGet("NATS_URL", "nats://nats:4222"),
		SMTP: SMTPConfig{
			Host: mustGet("SMTP_HOST", "smtp.gmail.com"),
			Port: mustGet("SMTP_PORT", "587"),
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
			From: os.Getenv("SMTP_FROM"),
		},
	}
}

func mustGet(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
