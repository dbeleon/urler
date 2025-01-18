package tnt

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/notifier/internal/domain/models"
	"github.com/dbeleon/urler/notifier/internal/queue"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	FuncPut     = "notif_put"
	FuncConsume = "notif_consume"
	FuncAck     = "notif_ack"
)

type TntQueue struct {
	conf Config
	conn *tarantool.Connection
}

func New(conf Config) *TntQueue {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  conf.Address,
		User:     conf.User,
		Password: conf.Password,
	}
	opts := tarantool.Opts{
		Timeout:       5 * time.Second,
		Reconnect:     conf.Reconnect,
		MaxReconnects: uint(conf.MaxReconnects),
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		log.Fatal("Connection refused:", zap.String("address", conf.Address), zap.Error(err))
	}

	return &TntQueue{
		conf: conf,
		conn: conn,
	}
}

func (t *TntQueue) Ack(id int64) error {
	var result []*AckResponse

	log.Debug("acking notification", zap.Int64("id", id))

	request := &AckRequest{
		Id: id,
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncAck).Args([]interface{}{request}))
	err := res.GetTyped(&result)
	if err != nil {
		return fmt.Errorf("ack notification task failed: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("unable to ack notification task: %w", err)
	}

	if result[0].Code != 0 {
		return fmt.Errorf("ack notification taskresponse error: code=%d, message=%s", result[0].Code, result[0].Message)
	}

	return nil
}

func (t *TntQueue) Consume() (*dom.NotifyTask, error) {
	var result []*ConsumeResponse

	request := &ConsumeRequest{
		Timeout: int(t.conf.Timeout),
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncConsume).Args([]interface{}{request}))
	err := res.GetTyped(&result)
	if err != nil {
		return nil, fmt.Errorf("consume notification task failed: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("unable to consume notification task: %w", err)
	}

	switch result[0].Code {
	case 1:
		return nil, queue.ErrEmptyQueue
	case 0:
		qr, err := base64.StdEncoding.DecodeString(result[0].QR)
		if err != nil {
			log.Error("decode qr code from task failed", zap.String("code_base64", result[0].QR), zap.String("url", result[0].Url))
			return nil, queue.ErrInvalidQRCode
		}

		return &dom.NotifyTask{Id: result[0].Id, Short: result[0].Url, UserIds: result[0].UserIDs, QR: qr}, nil
	default:
		return nil, fmt.Errorf("consume notification taskresponse error: code=%d, message=%s", result[0].Code, result[0].Message)
	}
}
