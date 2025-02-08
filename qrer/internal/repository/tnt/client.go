package tnt

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/qrer/internal/domain/models"
	"github.com/dbeleon/urler/qrer/internal/repository"
	mdl "github.com/dbeleon/urler/qrer/internal/repository/tnt/models"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
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
	confs    []Config
	connPool *pool.ConnectionPool
}

func New(confs []Config) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var instances []pool.Instance
	for _, conf := range confs {
		opts := tarantool.Opts{
			Timeout:       2 * time.Second,
			Reconnect:     conf.Reconnect,
			MaxReconnects: uint(conf.MaxReconnects),
		}
		log.Info("tnt instance", zap.String("name", conf.Address))
		inst := pool.Instance{
			Name: conf.Address,
			Dialer: tarantool.NetDialer{
				Address:  conf.Address,
				User:     conf.User,
				Password: conf.Password,
			},
			Opts: opts,
		}
		instances = append(instances, inst)
	}

	connPool, err := pool.Connect(ctx, instances)
	if err != nil {
		log.Fatal("Tnt connection refused:", zap.Error(err))
	}

	return &Client{
		confs:    confs,
		connPool: connPool,
	}
}

func (c *Client) Close() {
	c.connPool.CloseGraceful()
}

func (c *Client) SaveUrl(url dom.Url) (*dom.Url, error) {
	data := &mdl.Url{
		User:  url.User,
		Long:  url.Long,
		Short: url.Short,
	}

	res := c.connPool.Do(tarantool.NewCall17Request(FuncAddUrl).Args([]interface{}{data}), pool.RW)
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

func (c *Client) QRUpdate(qrTask dom.QRTask) ([]int64, error) {
	data := &mdl.Url{
		Short: qrTask.Short,
		QR:    base64.StdEncoding.EncodeToString(qrTask.QR),
	}

	res := c.connPool.Do(tarantool.NewCall17Request(FuncQRUpdate).Args([]interface{}{data}), pool.RW)
	var ans []*mdl.UserIDsResponse
	err := res.GetTyped(&ans)
	if err != nil {
		// var tntErr tarantool.Error
		// if errors.As(err, &tntErr) && tntErr.Code == iproto.ER_TUPLE_FOUND {
		// 	return nil, repository.ErrShortUrlCollision
		// }

		return nil, fmt.Errorf("qr code update failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	if ans[0].Code != 0 {
		return nil, fmt.Errorf("qr update failed: %w", errors.New(ans[0].Message))
	}

	return ans[0].UserIDs, nil
}

func (c *Client) GetUrl(short string) (*dom.Url, error) {
	data := &mdl.Url{
		Short: short,
	}
	res := c.connPool.Do(tarantool.NewCall17Request(FuncGetUrl).Args([]interface{}{data}), pool.RO)
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
