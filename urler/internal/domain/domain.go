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

type Config struct {
	Repo Repository
}

type Model struct {
	repo Repository
}

func New(conf Config) *Model {
	log.Debug("creating new model")
	return &Model{
		repo: conf.Repo,
	}
}

func (m *Model) MustStart() {
	log.Debug("starting model")
}

func (m *Model) Stop(ctx context.Context) {
	log.Debug("stopping model")
}
