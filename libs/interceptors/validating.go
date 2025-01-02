package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validator проверяет корректность данных
type Validator interface {
	Validate() error
}

// ValidatingInterceptor проверяет данные запросов с помощью интерфейса Validator
func ValidatingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	validator, ok := req.(Validator)
	if ok {
		err := validator.Validate()
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	return handler(ctx, req)
}
