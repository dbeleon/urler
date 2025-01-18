package repository

import "errors"

var (
	ErrNoResult          = errors.New("no result")
	ErrUserNotFound      = errors.New("user not found")
	ErrShortUrlCollision = errors.New("short url collision")
)
