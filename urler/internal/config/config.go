package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const defaultConfigPath = "/etc/urler/config.yaml"

type Config struct {
	Env             string `yaml:"env" env-default:"prod"`
	ShutdownTimeout int    `yaml:"shutdown_timeout" env-default:"5"`
	HTTPServer      `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:":8080"`
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
