package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbeleon/urler/libs/log"
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

	log.Init(cfg.Env == evnDevel)
	log.Info("app started")

	tntConf := tnt.Config{
		Address:       cfg.UrlsTntDB.Address,
		Reconnect:     time.Duration(cfg.UrlsTntDB.Reconnect) * time.Second,
		MaxReconnects: cfg.UrlsTntDB.MaxReconnects,
		User:          cfg.UrlsTntDB.User,
		Password:      cfg.UrlsTntDB.Password,
	}

	queueConf := queue.Config{
		Address:       cfg.QRTntQueue.Address,
		Reconnect:     time.Duration(cfg.QRTntQueue.Reconnect) * time.Second,
		MaxReconnects: cfg.QRTntQueue.MaxReconnects,
		User:          cfg.QRTntQueue.User,
		Password:      cfg.QRTntQueue.Password,
		Timeout:       cfg.QRTntQueue.Timeout,
	}

	conf := domain.Options{
		Repo:    tnt.New(tntConf),
		QRQueue: queue.New(queueConf),
		QR:      qr.New(),
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
