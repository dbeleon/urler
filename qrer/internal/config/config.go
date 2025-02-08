package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "/etc/qrer/config.yaml"

type Config struct {
	Env             string      `yaml:"env" env-default:"prod"`
	ShutdownTimeout int         `yaml:"shutdown_timeout" env-default:"5"`
	UrlTntDBs       []UrlsTntDB `yaml:"urls_tnts"`

	QRTntQueue `yaml:"qr_tnt"`
	Metrics    `yaml:"metrics"`
}

type UrlsTntDB struct {
	Address       string `yaml:"address"`
	Reconnect     int    `yaml:"reconnect"`
	MaxReconnects int    `yaml:"reconnects"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
}

type QRTntQueue struct {
	Address       string `yaml:"address"`
	Reconnect     int    `yaml:"reconnect"`
	MaxReconnects int    `yaml:"reconnects"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Timeout       uint   `yaml:"timeout"`
	Priority      uint   `yaml:"prior"`
	TTL           uint   `yaml:"ttl"`
	Delay         uint   `yaml:"delay"`
	TTR           uint   `yaml:"ttr"`
}

func MustLoad() *Config {
	configPath := os.Getenv("QRER_CONFIG")
	if configPath == "" {
		log.Print("QRER_CONFIG is not set, use default:", defaultConfigPath)
		configPath = defaultConfigPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read config file: %w", err))
	}

	log.Print(string(data))

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return &cfg
}

type Metrics struct {
	Address string `yaml:"address" env-default:":8888`
}
