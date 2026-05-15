package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway-service/internal/app"
	"gateway-service/internal/config"

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

	logger.Info("gateway started", zap.String("port", cfg.HTTPPort))

	go func() {
		if err := a.Run(); err != nil {
			logger.Error("http server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.Stop(ctx)
}
