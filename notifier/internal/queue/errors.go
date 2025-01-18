package queue

import "errors"

var (
	ErrEmptyQueue    = errors.New("empty queue")
	ErrInvalidQRCode = errors.New("invalid qr code")
)
