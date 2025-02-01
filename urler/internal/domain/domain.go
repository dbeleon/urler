package domain

import (
	"context"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/urler/internal/domain/models"
)

type Repository interface {
	AddUser(models.User) (*models.User, error)
	SaveUrl(url models.Url) (*models.Url, error)
	GetUrl(short string) (*models.Url, error)
	GetShorts(limit int64, offset int64) ([]string, error)
}

type Queue interface {
	Put(models.QRTask) (int, error)
}

type Config struct {
	Host string
}

type Options struct {
	Repo    Repository
	QRQueue Queue
}

type Model struct {
	conf    Config
	repo    Repository
	qrQueue Queue
}

func New(conf Config, opt Options) *Model {
	log.Debug("creating new model")
	return &Model{
		conf:    conf,
		repo:    opt.Repo,
		qrQueue: opt.QRQueue,
	}
}

func (m *Model) MustStart() {
	log.Debug("starting model")
}

func (m *Model) Stop(ctx context.Context) {
	log.Debug("stopping model")
}
