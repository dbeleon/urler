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
}

type Queue interface {
	Publish(models.QRTask) (int, error)
}

type Config struct {
	Host    string
	Repo    Repository
	QRQueue Queue
}

type Model struct {
	host    string
	repo    Repository
	qrQueue Queue
}

func New(conf Config) *Model {
	log.Debug("creating new model")
	return &Model{
		host:    conf.Host,
		repo:    conf.Repo,
		qrQueue: conf.QRQueue,
	}
}

func (m *Model) MustStart() {
	log.Debug("starting model")
}

func (m *Model) Stop(ctx context.Context) {
	log.Debug("stopping model")
}
