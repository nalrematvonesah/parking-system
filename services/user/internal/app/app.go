package app

import (
	"context"
	"net"

	"user-service/internal/config"
	"user-service/internal/email"
	"user-service/internal/handler"
	"user-service/internal/middleware"
	natspub "user-service/internal/nats"
	"user-service/internal/repository"
	"user-service/internal/service"

	userv1 "github.com/nalrematvonesah/user.proto/gen/user/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	grpcServer *grpc.Server
	lis        net.Listener
	publisher  *natspub.Publisher
}

func New(ctx context.Context, cfg config.Config, log *zap.Logger) (*App, error) {
	db, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	publisher, err := natspub.New(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	mailer := email.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.From)

	repo := repository.NewPostgres(db)
	svc := service.New(repo, publisher, mailer, log)
	h := handler.NewGRPC(svc)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.Recovery(log),
			middleware.UnaryLogger(log),
		),
	)

	userv1.RegisterUserServiceServer(server, h)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	return &App{grpcServer: server, lis: lis, publisher: publisher}, nil
}

func (a *App) Run() error {
	return a.grpcServer.Serve(a.lis)
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
	a.publisher.Close()
}
