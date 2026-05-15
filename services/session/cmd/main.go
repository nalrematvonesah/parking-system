package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"session-service/internal/app"
	"session-service/internal/config"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	a, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("failed to init app", zap.Error(err))
	}

	logger.Info("session-service started", zap.String("port", cfg.GRPCPort))

	go func() {
		if err := a.Run(); err != nil {
			logger.Error("grpc server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down gracefully...")
	a.Stop()
}
