package domain

import (
	"context"
	"github.com/dbeleon/urler/libs/log"
)

type Config struct {
}

type Model struct {
	config Config
}

func New(conf Config) *Model {
	log.Debug("creating new model")
	return &Model{
		config: conf,
	}
}

func (m *Model) MustStart() {
	log.Debug("starting model")
}

func (m *Model) Stop(ctx context.Context) {
	log.Debug("stopping model")
}
