package tnt

import (
	"context"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/qrer/internal/domain/models"
	"github.com/dbeleon/urler/qrer/internal/queue"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	FuncPublish = "qr_publish"
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

func (t *TntQueue) Publish(task dom.QRTask) (int, error) {
	log.Debug("publishing url", zap.String("short", task.Short))

	request := &PublishRequest{
		Url:      task.Short,
		Priority: task.Priority,
		TTL:      task.TTL,
		Delay:    task.Delay,
		TTR:      task.TTR,
	}

	if request.Priority == 0 {
		request.Priority = t.conf.Priority
	}

	if request.TTR == 0 {
		request.TTR = t.conf.TTR
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncPublish).Args([]interface{}{request}))
	var ans []*PublishResponse
	err := res.GetTyped(&ans)
	if err != nil {
		return 0, fmt.Errorf("publish qr task failed: %w", err)
	}

	if len(ans) == 0 {
		return 0, fmt.Errorf("unable to publish task: %w", err)
	}

	if ans[0].Code != 0 {
		return 0, fmt.Errorf("response error: code=%d, message=%s", ans[0].Code, ans[0].Message)
	}

	return int(ans[0].Id), nil
}

func (t *TntQueue) Ack(id int64) error {
	var result []*AckResponse

	log.Debug("acking", zap.Int64("id", id))

	request := &AckRequest{
		Id: id,
	}

	res := t.conn.Do(tarantool.NewCall17Request(FuncAck).Args([]interface{}{request}))
	err := res.GetTyped(&result)
	if err != nil {
		return fmt.Errorf("ack qr task failed: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("unable to ack task: %w", err)
	}

	if result[0].Code != 0 {
		return fmt.Errorf("response error: code=%d, message=%s", result[0].Code, result[0].Message)
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
		return nil, fmt.Errorf("response error: code=%d, message=%s", result[0].Code, result[0].Message)
	}
}

// func (t *tarqueue) pruneQueue(ctx context.Context) error {
// 	var result interface{}

// 	if err := t.client.Call17(ImporterPruneMethod, []interface{}{}, &result); err != nil {
// 		return err
// 	}

// 	t.logger.Ctx(ctx).Debugf("importer_prune res={%v}", result)

// 	return nil
// }
