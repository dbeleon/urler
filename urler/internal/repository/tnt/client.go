package tnt

import (
	"context"
	"errors"
	"fmt"
	"time"

	dom "github.com/dbeleon/urler/urler/internal/domain/models"
	"github.com/dbeleon/urler/urler/internal/repository"
	mdl "github.com/dbeleon/urler/urler/internal/repository/tnt/models"
	"go.uber.org/zap"

	"github.com/dbeleon/urler/libs/log"
	"github.com/tarantool/go-iproto"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
)

const (
	FuncAddUser   = "user_add"
	FuncAddUrl    = "url_add"
	FuncGetUrl    = "url_get"
	FuncGetShorts = "url_shorts"
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

func (c *Client) AddUser(usr dom.User) (*dom.User, error) {
	data := &mdl.User{
		Name:  usr.Name,
		Email: usr.Email,
	}
	res := c.connPool.Do(tarantool.NewCall17Request(FuncAddUser).Args([]interface{}{data}), pool.RW)
	var ans []*mdl.UserResponse
	err := res.GetTyped(&ans)
	if err != nil {
		return nil, fmt.Errorf("add user failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	if ans[0].Code != 0 {
		return nil, errors.New(ans[0].Message)
	}

	return &dom.User{
		Id:    ans[0].Id,
		Name:  usr.Name,
		Email: usr.Email,
	}, nil
}

func (c *Client) SaveUrl(url dom.Url) (*dom.Url, error) {
	data := &mdl.Url{
		User:  url.User,
		Long:  url.Long,
		Short: url.Short,
	}

	res := c.connPool.Do(tarantool.NewCall17Request(FuncAddUrl).Args([]interface{}{data}), pool.RW)
	var ans []*mdl.UrlResponse
	err := res.GetTyped(&ans)
	if err != nil {
		log.Error("tnt save url failed", zap.Error(err), zap.String("url", url.Long), zap.String("short", url.Short))
		var tntErr tarantool.Error
		if errors.As(err, &tntErr) && tntErr.Code == iproto.ER_TUPLE_FOUND {
			return nil, repository.ErrShortUrlCollision
		}

		return nil, fmt.Errorf("add url failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	switch ans[0].Code {
	case 0:
	case 2:
		return nil, repository.ErrUserNotFound
	default:
		return nil, errors.New(ans[0].Message)
	}

	url.Short = ans[0].Url.Short

	return &url, nil
}

func (c *Client) GetUrl(short string) (*dom.Url, error) {
	data := &mdl.Url{
		Short: short,
	}
	res := c.connPool.Do(tarantool.NewCall17Request(FuncGetUrl).Args([]interface{}{data}), pool.PreferRO)
	var ans []*mdl.UrlResponse
	err := res.GetTyped(&ans)
	if err != nil {
		return nil, fmt.Errorf("get url failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	if ans[0].Code != 0 {
		return nil, errors.New(ans[0].Message)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	url := &dom.Url{Short: short, Long: ans[0].Url.Long}

	return url, nil
}

func (c *Client) GetShorts(limit int64, offset int64) ([]string, error) {
	data := &mdl.LimOff{
		Limit:  limit,
		Offset: offset,
	}
	res := c.connPool.Do(tarantool.NewCall17Request(FuncGetShorts).Args([]interface{}{data}), pool.RO)
	var ans []*mdl.ShortsResponse
	err := res.GetTyped(&ans)
	if err != nil {
		return nil, fmt.Errorf("get shorts failed: %w", err)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	if ans[0].Code != 0 {
		return nil, errors.New(ans[0].Message)
	}

	if len(ans) == 0 {
		return nil, repository.ErrNoResult
	}

	return ans[0].Shorts, nil
}
