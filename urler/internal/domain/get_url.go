package domain

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUrlNotFound  = errors.New("url not found")
)

func (m *Model) GetUrl(ctx context.Context, user int64, shortUrl string) (string, error) {
	return "", ErrUserNotFound
}
