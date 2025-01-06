package domain

import (
	"context"
	"errors"
	"math/rand"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/urler/internal/domain/models"
	"github.com/dbeleon/urler/urler/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrNotAvailable    = errors.New("service is not available")
	ErrFailedToSaveUrl = errors.New("failed to save url")
)

func (m *Model) MakeUrl(ctx context.Context, user int64, long string) (string, error) {
	hashSize := 8
	url := models.Url{
		User:  user,
		Long:  long,
		Short: GenHash(hashSize),
	}
	res, err := m.repo.SaveUrl(url)
	for range 3 {
		if errors.Is(err, repository.ErrShortUrlCollision) {
			url.Short = GenHash(hashSize)
			res, err = m.repo.SaveUrl(url)
			continue
		}

		break
	}

	if err != nil {
		e := ErrNotRespond
		if errors.Is(err, repository.ErrShortUrlCollision) {
			e = ErrToManyHashColiisions
		}
		log.Error(e.Error(), zap.Error(err))
		return "", e
	}

	// TODO: use transactional outbox pattern
	if res.Short == url.Short {
		_, err = m.qrQueue.Publish(models.QRTask{
			Short: res.Short,
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
