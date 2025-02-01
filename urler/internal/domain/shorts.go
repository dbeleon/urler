package domain

import (
	"context"
	"errors"

	"github.com/dbeleon/urler/libs/log"
	"go.uber.org/zap"
)

var ErrShortsNotFound = errors.New("short urls not found")

func (m *Model) GetShorts(ctx context.Context, limit int64, offset int64) ([]string, error) {
	shorts, err := m.repo.GetShorts(limit, offset)
	if err != nil {
		log.Error(ErrShortsNotFound.Error(), zap.Error(err))
		return nil, ErrShortsNotFound
	}

	return shorts, nil
}
