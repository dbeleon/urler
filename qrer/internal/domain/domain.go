package domain

import (
	"context"
	"errors"
	"net/url"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/qrer/internal/domain/models"
	"github.com/dbeleon/urler/qrer/internal/queue"

	"go.uber.org/zap"
)

const workersNumber = 5

type Repository interface {
	SaveUrl(url models.Url) (*models.Url, error)
	QRUpdate(qrTask models.QRTask) ([]int64, error)
}

type Queuer interface {
	Put(models.NotifTask) error
	Consume() (*models.QRTask, error)
	Ack(int64) error
}

type QREncoder interface {
	Encode(text string) ([]byte, error)
}

type Options struct {
	Repo  Repository
	Queue Queuer
	QR    QREncoder
}

type Model struct {
	repo    Repository
	qrQueue Queuer
	qr      QREncoder

	doneChan chan struct{}
}

func New(opt Options) *Model {
	log.Debug("creating new model")
	return &Model{
		repo:    opt.Repo,
		qrQueue: opt.Queue,
		qr:      opt.QR,

		doneChan: make(chan struct{}),
	}
}

func (m *Model) MustStart() {
	log.Debug("starting model")
	for i := 0; i < workersNumber; i++ {
		go m.Worker()
	}
}

func (m *Model) Stop(ctx context.Context) {
	log.Debug("stopping model")
	close(m.doneChan)
}

func (m *Model) Worker() {
	for {
		select {
		case <-m.doneChan:
			return
		default:
		}

		task, err := m.qrQueue.Consume()
		if errors.Is(err, queue.ErrEmptyQueue) {
			continue
		}
		if err != nil {
			log.Error("qr queue consume task failed", zap.Error(err))
			continue
		}

		addr, err := url.JoinPath(task.Host, task.Short)
		if err != nil {
			log.Error("cannot join path", zap.Error(err))
			continue
		}

		qrData, err := m.qr.Encode(addr)
		if err != nil {
			log.Error("QRCode generation failed", zap.String("url", addr), zap.Error(err))
			continue
		}

		task.QR = qrData
		userIDs, err := m.repo.QRUpdate(*task)
		if err != nil {
			log.Error("could not update url qr", zap.Error(err))
			continue
		}

		m.qrQueue.Put(models.NotifTask{
			Short:   task.Short,
			UserIDs: userIDs,
			QR:      qrData,
		})

		err = m.qrQueue.Ack(task.Id)
		if err != nil {
			log.Error("qr queue ack task failed", zap.Int64("qr_task_id", task.Id), zap.Error(err))
			continue
		}
	}
}
