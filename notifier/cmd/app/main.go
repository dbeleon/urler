package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/notifier/internal/config"
	"github.com/dbeleon/urler/notifier/internal/domain"
	queue "github.com/dbeleon/urler/notifier/internal/queue/tnt"
)

const (
	evnDevel = "devel"
	evnProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log.Init(cfg.Env == evnDevel)
	log.Info("app started")

	queueConf := queue.Config{
		Address:       cfg.TntQueue.Address,
		Reconnect:     time.Duration(cfg.TntQueue.Reconnect) * time.Second,
		MaxReconnects: cfg.TntQueue.MaxReconnects,
		User:          cfg.TntQueue.User,
		Password:      cfg.TntQueue.Password,
		Timeout:       cfg.TntQueue.Timeout,
	}

	conf := domain.Options{
		Queue: queue.New(queueConf),
	}

	app := domain.New(conf)

	app.MustStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctxOut, shutdown := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer shutdown()

	log.Info("app exiting gracefully")

	app.Stop(ctxOut)
}
