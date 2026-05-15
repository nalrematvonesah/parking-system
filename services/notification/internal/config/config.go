package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	NatsURL string
	SMTP    SMTPConfig
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
		NatsURL: get("NATS_URL", "nats://nats:4222"),
		SMTP: SMTPConfig{
			Host: get("SMTP_HOST", "smtp.gmail.com"),
			Port: get("SMTP_PORT", "587"),
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
			From: os.Getenv("SMTP_FROM"),
		},
	}
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
