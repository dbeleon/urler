package domain

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/notifier/internal/domain/models"
	"github.com/dbeleon/urler/notifier/internal/queue"

	"go.uber.org/zap"
)

const workersNumber = 12

type NotifyQueue interface {
	Consume() (*models.NotifyTask, error)
	Ack(int64) error
}

type Options struct {
	Queue NotifyQueue
}

type Model struct {
	queue NotifyQueue

	doneChan chan struct{}
}

func New(opt Options) *Model {
	log.Debug("creating new model")
	return &Model{
		queue: opt.Queue,

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
	ack := func(id int64) {
		err := m.queue.Ack(id)
		if err != nil {
			log.Error("notification queue ack task failed", zap.Int64("notif_task_id", id), zap.Error(err))
		}
	}

	for {
		select {
		case <-m.doneChan:
			return
		default:
		}

		task, err := m.queue.Consume()
		if errors.Is(err, queue.ErrEmptyQueue) {
			continue
		}
		if err != nil {
			log.Error("queue consume notification task failed", zap.Error(err))
			if errors.Is(err, queue.ErrInvalidQRCode) {
				ack(task.Id)
			}
			continue
		}

		log.Info("url qr code created",
			zap.String("url", task.Short),
			zap.String("users", fmt.Sprint(task.UserIds)),
			zap.String("code_base64", base64.StdEncoding.EncodeToString(task.QR)))

		err = m.queue.Ack(task.Id)
		if err != nil {
			log.Error("notification queue ack task failed", zap.Int64("notif_task_id", task.Id), zap.Error(err))
			continue
		}
	}
}
