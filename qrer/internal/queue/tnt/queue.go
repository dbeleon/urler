package tnt

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/qrer/internal/domain/models"
	"github.com/dbeleon/urler/qrer/internal/queue"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	FuncPut     = "notif_put"
	FuncConsume = "qr_consume"
	FuncAck     = "qr_ack"
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

func (t *TntQueue) Put(task dom.NotifTask) error {
	log.Debug("putting notification", zap.String("short", task.Short))

	request := &PutRequest{
		Url:     task.Short,
		UserIDs: task.UserIDs,
		QR:      base64.StdEncoding.EncodeToString(task.QR),
	}

	if request.Priority == 0 {
		request.Priority = t.conf.Priority
	}

	if request.TTR == 0 {
		request.TTR = t.conf.TTR
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncPut).Args([]interface{}{request}))
	var ans []*PutResponse
	err := res.GetTyped(&ans)
	if err != nil {
		return fmt.Errorf("put notification task failed: %w", err)
	}

	if len(ans) == 0 {
		return fmt.Errorf("unable to put notification task: %w", err)
	}

	if ans[0].Code != 0 {
		return fmt.Errorf("put notification response error: code=%d, message=%s", ans[0].Code, ans[0].Message)
	}

	return nil
}

func (t *TntQueue) Ack(id int64) error {
	var result []*AckResponse

	log.Debug("acking qr", zap.Int64("id", id))

	request := &AckRequest{
		Id: id,
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncAck).Args([]interface{}{request}))
	err := res.GetTyped(&result)
	if err != nil {
		return fmt.Errorf("ack qr task failed: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("unable to ack qr task: %w", err)
	}

	if result[0].Code != 0 {
		return fmt.Errorf("qr task ack response error: code=%d, message=%s", result[0].Code, result[0].Message)
	}

	return nil
}

func (t *TntQueue) Consume() (*dom.QRTask, error) {
	var result []*ConsumeResponse

	request := &ConsumeRequest{
		Timeout: int(t.conf.Timeout),
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncConsume).Args([]interface{}{request}))
	err := res.GetTyped(&result)
	if err != nil {
		return nil, fmt.Errorf("consume qr task failed: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("unable to consume qr task: %w", err)
	}

	switch result[0].Code {
	case 1:
		return nil, queue.ErrEmptyQueue
	case 0:
		return &dom.QRTask{Id: result[0].Id, Short: result[0].Url}, nil
	default:
		return nil, fmt.Errorf("consume qr task response error: code=%d, message=%s", result[0].Code, result[0].Message)
	}
}
