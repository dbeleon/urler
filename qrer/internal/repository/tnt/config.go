package tnt

import "time"

type Config struct {
	Address       string
	Reconnect     time.Duration
	MaxReconnects int
	User          string
	Password      string
}
