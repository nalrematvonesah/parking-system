package main

import (
	"context"
	"log"

	"user-service/internal/app"
	"user-service/internal/config"

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

	logger.Info("user-service started", zap.String("port", cfg.GRPCPort))

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
