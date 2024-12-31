package main

import (
	"context"
	"github.com/dbeleon/urler/internal/domain"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbeleon/urler/internal/config"
	"github.com/dbeleon/urler/libs/log"
)

const (
	evnDevel = "devel"
	evnProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log.Init(cfg.Env == evnDevel)
	log.Info("app started")

	conf := domain.Config{}

	app := domain.New(conf)

	app.MustStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer shutdown()

	log.Info("app exiting gracefully")

	app.Stop(ctx)
}
