package app

import (
	"context"
	"net"

	"parking-service/internal/cache"
	"parking-service/internal/config"
	"parking-service/internal/handler"
	"parking-service/internal/middleware"
	"parking-service/internal/repository"
	"parking-service/internal/service"

	parkingpb "github.com/nalrematvonesah/parking-proto/gen/parking/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	lis        net.Listener
}

func New(ctx context.Context, cfg config.Config, log *zap.Logger) (*App, error) {
	db, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}

	repo := repository.NewPostgres(db)
	cache := cache.New(cfg.RedisAddr, cfg.RedisDB)
	svc := service.New(repo, cache)
	h := handler.NewGRPC(svc)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Recovery(log),
			middleware.UnaryLogger(log),
		),
	)

	parkingpb.RegisterParkingServiceServer(server, h)
	reflection.Register(server)
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &App{grpcServer: server, lis: lis}, nil
}

func (a *App) Run() error {
	return a.grpcServer.Serve(a.lis)
}
