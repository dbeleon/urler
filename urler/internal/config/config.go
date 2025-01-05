package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "/etc/urler/config.yaml"

type Config struct {
	Env             string `yaml:"env" env-default:"prod"`
	ShutdownTimeout int    `yaml:"shutdown_timeout" env-default:"5"`
	GRPCServer      `yaml:"grpc_server"`
	HTTPServer      `yaml:"http_server"`
	UrlsTntDB       `yaml:"urls_tnt"`
	QRTntQueue      `yaml:"qr_tnt"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:":8000"`
}

type GRPCServer struct {
	Address string `yaml:"address" env-default:":8080"`
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
}

func MustLoad() *Config {
	configPath := os.Getenv("URLER_CONFIG")
	if configPath == "" {
		log.Print("URLER_CONFIG is not set, use default:", defaultConfigPath)
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
