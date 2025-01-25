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
)

const (
	FuncAddUser = "user_add"
	FuncAddUrl  = "url_add"
	FuncGetUrl  = "url_get"
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

func (c *Client) AddUser(usr dom.User) (*dom.User, error) {
	data := &mdl.User{
		Name:  usr.Name,
		Email: usr.Email,
	}
	res := c.conn.Do(tarantool.NewCall17Request(FuncAddUser).Args([]interface{}{data}))
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

	res := c.conn.Do(tarantool.NewCall17Request(FuncAddUrl).Args([]interface{}{data}))
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
	res := c.conn.Do(tarantool.NewCall17Request(FuncGetUrl).Args([]interface{}{data}))
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
