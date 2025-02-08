package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/libs/metrics"
	"github.com/dbeleon/urler/qrer/internal/config"
	"github.com/dbeleon/urler/qrer/internal/domain"
	"github.com/dbeleon/urler/qrer/internal/qr"
	queue "github.com/dbeleon/urler/qrer/internal/queue/tnt"
	"github.com/dbeleon/urler/qrer/internal/repository/tnt"
)

const (
	evnDevel = "devel"
	evnProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serviceMetrics := metrics.New(cfg.Metrics.Address)
	defer serviceMetrics.Close(ctx)

	log.Init(cfg.Env == evnDevel)
	log.Info("app started")

	tntConfs := make([]tnt.Config, 0, len(cfg.UrlTntDBs))
	for _, c := range cfg.UrlTntDBs {
		tntConfs = append(tntConfs, tnt.Config{
			Address:       c.Address,
			Reconnect:     time.Duration(c.Reconnect) * time.Second,
			MaxReconnects: c.MaxReconnects,
			User:          c.User,
			Password:      c.Password,
		})
	}
	tntClient := tnt.New(tntConfs)

	queueConf := queue.Config{
		Address:       cfg.QRTntQueue.Address,
		Reconnect:     time.Duration(cfg.QRTntQueue.Reconnect) * time.Second,
		MaxReconnects: cfg.QRTntQueue.MaxReconnects,
		User:          cfg.QRTntQueue.User,
		Password:      cfg.QRTntQueue.Password,
		Timeout:       cfg.QRTntQueue.Timeout,
	}

	conf := domain.Options{
		Repo:  tntClient,
		Queue: queue.New(queueConf),
		QR:    qr.New(),
	}

	app := domain.New(conf)

	app.MustStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctxOut, shutdown := context.WithTimeout(ctx, time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer shutdown()

	log.Info("app exiting gracefully")

	app.Stop(ctxOut)
}
