package tnt

import (
	"context"
	"errors"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/qrer/internal/domain/models"
	mdl "github.com/dbeleon/urler/qrer/internal/repository/tnt/models"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	FuncAddUser  = "user_add"
	FuncAddUrl   = "url_add"
	FuncGetUrl   = "url_get"
	FuncQRUpdate = "qr_update"
)

// Client connects to tarantool
// TODO: retry, models convert
type Client struct {
	conf Config
	conn *tarantool.Connection
}

func New(conf Config) *Client {
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

	return &Client{
		conf: conf,
		conn: conn,
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) SaveUrl(url dom.Url) (*dom.Url, error) {
	data := &mdl.Url{
		User:  url.User,
		Long:  url.Long,
		Short: url.Short,
	}

	res := c.conn.Do(tarantool.NewCall17Request(FuncAddUrl).Args([]interface{}{data}))
	var ans []*mdl.Url
	err := res.GetTyped(&ans)
	if err != nil {
		// var tntErr tarantool.Error
		// if errors.As(err, &tntErr) && tntErr.Code == iproto.ER_TUPLE_FOUND {
		// 	return nil, repository.ErrShortUrlCollision
		// }

		return nil, fmt.Errorf("add url failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, fmt.Errorf("unable to add url: %w", err)
	}

	url.Short = ans[0].Short

	return &url, nil
}

func (c *Client) QRUpdate(qrTask dom.QRTask) error {
	data := &mdl.Url{
		Short: qrTask.Short,
		QR:    qrTask.QR,
	}

	res := c.conn.Do(tarantool.NewCall17Request(FuncQRUpdate).Args([]interface{}{data}))
	var ans []*mdl.BaseResponse
	err := res.GetTyped(&ans)
	if err != nil {
		// var tntErr tarantool.Error
		// if errors.As(err, &tntErr) && tntErr.Code == iproto.ER_TUPLE_FOUND {
		// 	return nil, repository.ErrShortUrlCollision
		// }

		return fmt.Errorf("qr update failed: %w", err)
	}

	if len(ans) == 0 {
		return fmt.Errorf("unable to update qr: %w", err)
	}

	switch ans[0].Code {
	case 0:
		return nil
	default:
	}

	return fmt.Errorf("qr update failed: %w", errors.New(ans[0].Message))
}

func (c *Client) GetUrl(short string) (*dom.Url, error) {
	data := &mdl.Url{
		Short: short,
	}
	res := c.conn.Do(tarantool.NewCall17Request(FuncGetUrl).Args([]interface{}{data}))
	var ans []*mdl.Url
	err := res.GetTyped(&ans)
	if err != nil {
		return nil, fmt.Errorf("get url failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, fmt.Errorf("unable to get url: %w", err)
	}

	url := &dom.Url{Short: short, Long: ans[0].Long}

	return url, nil
}
