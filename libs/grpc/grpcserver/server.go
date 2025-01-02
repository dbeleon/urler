package grpcserver

import (
	"net"

	"github.com/dbeleon/urler/libs/interceptors"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server это gRPC сервер с метриками и трейсами
type server struct {
	srv *grpc.Server
}

// New создает новый server
func New() *server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptors.LoggingInterceptor,
				interceptors.ValidatingInterceptor,
			),
		),
	)

	reflection.Register(s)

	return &server{
		srv: s,
	}
}

// RegisterService registers a service and its implementation to the
// concrete type implementing this interface.  It may not be called
// once the server has started serving.
// desc describes the service and its methods and handlers. impl is the
// service implementation which is passed to the method handlers.
func (s *server) RegisterService(desc *grpc.ServiceDesc, ss interface{}) {
	s.srv.RegisterService(desc, ss)
}

// GetServiceInfo returns a map from service names to ServiceInfo.
// Service names include the package names, in the form of <package>.<service>.
func (s *server) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.srv.GetServiceInfo()
}

// Serve accepts incoming connections on the listener lis, creating a new
// ServerTransport and service goroutine for each. The service goroutines
// read gRPC requests and then call the registered handlers to reply to them.
// Serve returns when lis.Accept fails with fatal errors.  lis will be closed when
// this method returns.
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s *server) Serve(lis net.Listener) error {
	return s.srv.Serve(lis)
}
