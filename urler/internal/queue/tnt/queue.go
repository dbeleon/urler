package tnt

import (
	"context"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/urler/internal/domain/models"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	FuncPut = "qr_put"
)

type tarqueue struct {
	conf Config
	conn *tarantool.Connection
}

func New(conf Config) *tarqueue {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  conf.Address,
		User:     conf.User,
		Password: conf.Password,
	}
	opts := tarantool.Opts{
		Timeout:       2 * time.Second,
		Reconnect:     conf.Reconnect,
		MaxReconnects: uint(conf.MaxReconnects),
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		log.Fatal("Connection refused:", zap.String("address", conf.Address), zap.Error(err))
	}

	return &tarqueue{
		conf: conf,
		conn: conn,
	}
}

func (t *tarqueue) Put(task dom.QRTask) (int, error) {
	log.Debug("putting url", zap.String("short", task.Short))

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

	res := t.conn.Do(tarantool.NewCall17Request(FuncPut).Args([]interface{}{request}))
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
