package main

import (
	"context"
	"log"

	"parking-service/internal/app"
	"parking-service/internal/config"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()
	logger, _ := zap.NewProduction()

	a, err := app.New(ctx, cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("parking-service started")

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
