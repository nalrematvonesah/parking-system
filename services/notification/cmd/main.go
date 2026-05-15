package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/app"
	"notification-service/internal/config"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Sync()

	a, err := app.New(cfg, logger)
	if err != nil {
		logger.Fatal("init app", zap.Error(err))
	}

	logger.Info("notification-service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down gracefully...")
	a.Stop()
}
