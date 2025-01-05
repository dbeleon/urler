package domain

import (
	"context"
	"errors"

	"github.com/dbeleon/urler/libs/log"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrShortUrlNotFound     = errors.New("short url not found")
	ErrNotRespond           = errors.New("repository is not responding")
	ErrToManyHashColiisions = errors.New("to many url hash collisions")
)

func (m *Model) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	url, err := m.repo.GetUrl(shortUrl)
	if err != nil {
		log.Error(ErrShortUrlNotFound.Error(), zap.Error(err))
		return "", ErrShortUrlNotFound
	}

	return url.Long, nil
}
