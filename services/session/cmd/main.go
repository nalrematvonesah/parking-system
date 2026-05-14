package main

import (
	"context"
	"log"

	"session-service/internal/app"
	"session-service/internal/cache"
	"session-service/internal/config"
	"session-service/internal/handler"
	"session-service/internal/nats"
	"session-service/internal/repository"
	"session-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	db, err := pgxpool.New(
		context.Background(),
		cfg.DBUrl,
	)

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)

	parkingClient, err := service.NewParkingClient(
		"localhost:50051",
	)

	if err != nil {
		log.Fatal(err)
	}

	redisCache := cache.New("localhost:6379")

	natsClient, err := nats.New(
		"nats://localhost:4222",
	)

	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(
		repo,
		parkingClient,
		redisCache,
		natsClient,
	)

	h := handler.New(svc)

	application := app.New(
		cfg,
		h,
	)

	log.Fatal(application.Run())
}
