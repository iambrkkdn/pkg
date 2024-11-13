package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func GRPCLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Выполняем обработку запроса
		res, err := handler(ctx, req)

		// Логгируем информацию о запросе после завершения
		duration := time.Since(start)
		clientIP := ""
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		logger.Info("gRPC Request",
			zap.String("method", info.FullMethod),
			zap.String("client_ip", clientIP),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return res, err
	}
}
