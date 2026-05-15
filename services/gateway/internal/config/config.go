package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string

	UserAddr    string
	ParkingAddr string
	SessionAddr string

	JWTSecret string
	JWTTTL    time.Duration

	AllowedOrigins []string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	return Config{
		HTTPPort:       get("HTTP_PORT", "8080"),
		UserAddr:       get("USER_ADDR", "user:50052"),
		ParkingAddr:    get("PARKING_ADDR", "parking:50051"),
		SessionAddr:    get("SESSION_ADDR", "session:50053"),
		JWTSecret:      get("JWT_SECRET", "dev-secret-change-me"),
		JWTTTL:         24 * time.Hour,
		AllowedOrigins: parseOrigins(get("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
	}
}

func parseOrigins(s string) []string {
	var out []string
	for _, o := range strings.Split(s, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			out = append(out, o)
		}
	}
	return out
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
