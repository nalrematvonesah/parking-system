package app

import (
	"context"
	"net"

	"session-service/internal/cache"
	"session-service/internal/config"
	"session-service/internal/handler"
	"session-service/internal/middleware"
	natspub "session-service/internal/nats"
	"session-service/internal/repository"
	"session-service/internal/service"

	sessionv1 "github.com/nalrematvonesah/session.proto/gen/session/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer    *grpc.Server
	lis           net.Listener
	publisher     *natspub.Publisher
	parkingClient *service.ParkingClient
	db            *pgxpool.Pool
}

func New(ctx context.Context, cfg config.Config, log *zap.Logger) (*App, error) {
	db, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	parkingClient, err := service.NewParkingClient(cfg.ParkingAddr)
	if err != nil {
		return nil, err
	}

	publisher, err := natspub.New(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	repo := repository.New(db)
	redisCache := cache.New(cfg.RedisAddr)

	svc := service.New(repo, parkingClient, redisCache, publisher, log, cfg.PricePerHour)
	h := handler.New(svc)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Recovery(log),
			middleware.UnaryLogger(log),
		),
	)
	sessionv1.RegisterSessionServiceServer(server, h)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &App{
		grpcServer:    server,
		lis:           lis,
		publisher:     publisher,
		parkingClient: parkingClient,
		db:            db,
	}, nil
}

func (a *App) Run() error {
	return a.grpcServer.Serve(a.lis)
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
	a.publisher.Close()
	_ = a.parkingClient.Close()
	a.db.Close()
}
