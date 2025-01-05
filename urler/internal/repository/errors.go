package repository

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrShortUrlCollision = errors.New("short url collision")
)
