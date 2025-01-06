package domain

import (
	"bytes"
	"context"
	"errors"
	"net/url"

	"github.com/dbeleon/urler/libs/log"
	"github.com/dbeleon/urler/qrer/internal/domain/models"
	"github.com/dbeleon/urler/qrer/internal/queue"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"go.uber.org/zap"
)

type Repository interface {
	SaveUrl(url models.Url) (*models.Url, error)
	QRUpdate(qrTask models.QRTask) error
}

type Queue interface {
	Publish(models.QRTask) (int, error)
	Consume() (*models.QRTask, error)
	Ack(int64) error
}

type Config struct {
	Repo    Repository
	QRQueue Queue
}

type Model struct {
	repo    Repository
	qrQueue Queue

	doneChan chan struct{}
}

func New(conf Config) *Model {
	log.Debug("creating new model")
	return &Model{
		repo:     conf.Repo,
		qrQueue:  conf.QRQueue,
		doneChan: make(chan struct{}),
	}
}

func (m *Model) MustStart() {
	log.Debug("starting model")
	for i := 0; i < 1; i++ {
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
		qrc, err := qrcode.New(addr)
		if err != nil {
			log.Error("could not generate QRCode", zap.String("url", addr), zap.Error(err))
			continue
		}

		data := make([]byte, 0, 100*1024)
		buf := bytes.NewBuffer(data)
		wc := &WrCl{Buf: buf}
		w := standard.NewWithWriter(wc)
		if err = qrc.Save(w); err != nil {
			log.Error("could not save image", zap.Error(err))
			continue
		}
		task.QR = wc.Buf.Bytes()
		err = m.repo.QRUpdate(*task)
		if err != nil {
			log.Error("could not update url qr", zap.Error(err))
			//continue
		}

		// option := compressed.Option{
		// 	Padding:   4, // padding pixels around the qr code.
		// 	BlockSize: 1, // block pixels which represents a bit data.
		// }

		// writer, err := compressed.New("../assets/qrcode_small.png", &option)
		// if err != nil {
		// 	log.Error("QR code compressed.New write failed", zap.Error(err))
		// 	continue
		// }

		// if err := qrc.Save(writer); err != nil {
		// 	log.Error("could not save image", zap.Error(err))
		// 	continue
		// }

		err = m.qrQueue.Ack(task.Id)
		if err != nil {
			log.Error("qr queue ack task failed", zap.Int64("qr_task_id", task.Id), zap.Error(err))
			continue
		}
	}
}

type WrCl struct {
	Buf *bytes.Buffer
}

func (w *WrCl) Write(data []byte) (int, error) {
	return w.Buf.Write(data)
}

func (w *WrCl) Close() error {
	return nil
}
