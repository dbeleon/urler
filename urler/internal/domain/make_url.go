package domain

import (
	"context"
	"errors"
	"math/rand"

	"github.com/avast/retry-go/v4"
	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/urler/internal/domain/models"
	"github.com/dbeleon/urler/urler/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrNotAvailable    = errors.New("service is not available")
	ErrFailedToSaveUrl = errors.New("failed to save url")
)

const (
	SaveUrlAttempts   = 5
	CollisionAttempts = 5
)

func (m *Model) MakeUrl(ctx context.Context, user int64, long string) (string, error) {
	hashSize := 8
	urlModel := models.Url{
		User:  user,
		Long:  long,
		Short: GenHash(hashSize),
	}
	collisionAttempts := CollisionAttempts
	var res *models.Url
	err := retry.Do(
		func() error {
			var err error
			res, err = m.repo.SaveUrl(urlModel)
			for collisionAttempts > 0 {
				if errors.Is(err, repository.ErrShortUrlCollision) {
					collisionAttempts--
					urlModel.Short = GenHash(hashSize)
					res, err = m.repo.SaveUrl(urlModel)
					continue
				}

				break
			}

			return err
		},
		retry.Attempts(SaveUrlAttempts),
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

func GenHash(size int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
