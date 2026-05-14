package app

import (
	"log"
	"net"

	"session-service/internal/config"
	"session-service/internal/handler"

	"google.golang.org/grpc"
)

type App struct {
	cfg     *config.Config
	handler *handler.GRPCHandler
}

func New(
	cfg *config.Config,
	handler *handler.GRPCHandler,
) *App {
	return &App{
		cfg:     cfg,
		handler: handler,
	}
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.cfg.GRPCPort)
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	log.Println("session grpc running on", a.cfg.GRPCPort)

	return server.Serve(lis)
}
