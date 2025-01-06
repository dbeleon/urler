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
	urlModel := models.Url{
		User:  user,
		Long:  long,
		Short: GenHash(hashSize),
	}
	res, err := m.repo.SaveUrl(urlModel)
	for range 3 {
		if errors.Is(err, repository.ErrShortUrlCollision) {
			urlModel.Short = GenHash(hashSize)
			res, err = m.repo.SaveUrl(urlModel)
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

	// TODO: use transactional outbox
	if res.Short == urlModel.Short {
		if err != nil {
			log.Error("url path join failed", zap.Error(err))
		} else {
			_, err = m.qrQueue.Publish(models.QRTask{
				Host:  m.host,
				Short: res.Short,
				TTR:   10,
			})
			if err != nil {
				log.Error("publish to qr queue failed", zap.Error(err))
			}
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
