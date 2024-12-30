package domain

import (
	"context"
	"errors"
)

var (
	ErrNotAvailable = errors.New("service is not available")
)

func (m *Model) MakeUrl(ctx context.Context, user int64, url string) (string, error) {
	return "", ErrUserNotFound
}
