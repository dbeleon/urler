package domain

import (
	"context"
	"errors"

	"github.com/avast/retry-go/v4"
	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/urler/internal/domain/models"
	"github.com/dbeleon/urler/urler/internal/repository"
	"github.com/dbeleon/urler/urler/internal/tiny"
	"go.uber.org/zap"
)

var (
	ErrNotAvailable    = errors.New("service is not available")
	ErrFailedToSaveUrl = errors.New("failed to save url")
)

const (
	HashSize = 8
)

func (m *Model) MakeUrl(ctx context.Context, user int64, long string) (string, error) {
	hash := tiny.Get(long)
	offset := 0
	urlModel := models.Url{
		User:  user,
		Long:  long,
		Short: hash[offset:HashSize],
	}
	collisionAttempts := uint(len(hash) + 1 - HashSize)
	var res *models.Url
	err := retry.Do(
		func() error {
			var err error
			urlModel.Short = hash[offset:HashSize]
			res, err = m.repo.SaveUrl(urlModel)
			if errors.Is(err, repository.ErrShortUrlCollision) {
				offset++
			}

			return err
		},
		retry.Attempts(collisionAttempts),
		retry.OnRetry(func(n uint, err error) {
			log.Info("retrying to save url", zap.Int("attempt", int(n)), zap.Error(err),
				zap.String("long", urlModel.Long), zap.String("short", urlModel.Short))
		}),
	)

	if err != nil {
		e := ErrNotRespond
		if errors.Is(err, repository.ErrShortUrlCollision) {
			e = ErrToManyHashCollisions
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			e = ErrUserNotFound
		}
		log.Error(e.Error(), zap.Error(err))
		return "", e
	}

	// TODO: use transactional outbox - m.b. not need?
	if res.Short == urlModel.Short {
		_, err = m.qrQueue.Put(models.QRTask{
			Host:  m.conf.Host,
			Short: res.Short,
			TTR:   10,
		})
		if err != nil {
			log.Error("publish to qr queue failed", zap.Error(err))
		}
	}

	return res.Short, nil
}
