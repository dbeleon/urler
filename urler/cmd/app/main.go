package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/dbeleon/urler/libs/grpc/grpcserver"
	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/urler/internal/config"
	"github.com/dbeleon/urler/urler/internal/domain"
	"github.com/dbeleon/urler/urler/internal/repository/tnt"
	svc "github.com/dbeleon/urler/urler/internal/service/urler"
	urler "github.com/dbeleon/urler/urler/pkg/urler/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/zap"
)

const (
	evnDevel = "devel"
	evnProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log.Init(cfg.Env == evnDevel)
	log.Info("app started")

	tntConf := tnt.Config{
		Address:       cfg.Tnt.Address,
		Reconnect:     time.Duration(cfg.Tnt.Reconnect) * time.Second,
		MaxReconnects: cfg.Tnt.MaxReconnects,
		User:          cfg.Tnt.User,
		Password:      cfg.Tnt.Password,
	}
	tntClient := tnt.New(tntConf)

	conf := domain.Config{
		Repo: tntClient,
	}

	app := domain.New(conf)

	app.MustStart()

	lis, err := net.Listen("tcp", cfg.GRPCServer.Address)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err), zap.String("port", cfg.GRPCServer.Address))
	}

	s := grpcserver.New()
	urler.RegisterUrlerServiceServer(s, svc.New(app))

	log.Info("gRPC server listening", zap.String("address", lis.Addr().String()))

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithForwardResponseOption(svc.ResponseHeaderMatcher))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = urler.RegisterUrlerServiceHandlerFromEndpoint(ctx, mux, cfg.GRPCServer.Address, opts)
	if err != nil {
		log.Fatal("failed to register service handler from endpoint", zap.Error(err))
	}

	srv := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	log.Info("HTTP server listening", zap.String("address", srv.Addr))
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	cancel()

	ctxOut, shutdown := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer shutdown()

	log.Info("app exiting gracefully")

	s.GracefulStop()

	app.Stop(ctxOut)
}
