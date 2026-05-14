package config

type Config struct {
	DBUrl    string
	GRPCPort string
}

func Load() *Config {
	return &Config{
		DBUrl:    "postgres://postgres:123@localhost:55434/parking?sslmode=disable",
		GRPCPort: ":50052",
	}
}
