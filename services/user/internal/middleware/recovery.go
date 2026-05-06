package middleware

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Recovery(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered", zap.Any("panic", r))
			}
		}()
		return handler(ctx, req)
	}
}
