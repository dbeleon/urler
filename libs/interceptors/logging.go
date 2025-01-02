package interceptors

import (
	"context"

	"github.com/dbeleon/urler/libs/log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// LoggingInterceptor логирует вызовы gRPC методов
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Info("handle call", zap.String("module", "gRPC"), zap.String("handle", info.FullMethod), zap.Any("request", req))

	res, err := handler(ctx, req)
	if err != nil {
		log.Error("handle error", zap.Error(err), zap.String("module", "gRPC"), zap.String("handle", info.FullMethod), zap.Any("request", req))
		return nil, err
	}

	log.Info("handle result", zap.String("module", "gRPC"), zap.String("handle", info.FullMethod), zap.Any("result", res))

	return res, nil
}
