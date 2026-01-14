package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env         string `yaml:"env"`
	GRPCAddr    string `yaml:"grpc_addr"`
	WSAddr      string `yaml:"ws_addr"`
	MySQLDSN    string `yaml:"mysql_dsn"`
	RabbitMQURL string `yaml:"rabbitmq_url"`
}

func Load() Config {
	cfg := Config{
		Env:         getString("APP_ENV", "dev"),
		GRPCAddr:    getString("GRPC_ADDR", ":50051"),
		WSAddr:      getString("WS_ADDR", ":50052"),
		MySQLDSN:    getString("MYSQL_DSN", ""),
		RabbitMQURL: getString("RABBITMQ_URL", ""),
	}

	if path := getString("CONFIG_FILE", ""); path != "" {
		_ = loadFromYAML(path, &cfg)
	} else {
		_ = loadFromYAML("config.local.yaml", &cfg)
	}

	cfg.Env = getString("APP_ENV", cfg.Env)
	cfg.GRPCAddr = getString("GRPC_ADDR", cfg.GRPCAddr)
	cfg.WSAddr = getString("WS_ADDR", cfg.WSAddr)
	cfg.MySQLDSN = getString("MYSQL_DSN", cfg.MySQLDSN)
	cfg.RabbitMQURL = getString("RABBITMQ_URL", cfg.RabbitMQURL)
	return cfg
}

func getString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func loadFromYAML(path string, out any) error {
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read config: %w", err)
	}
	if err := yaml.Unmarshal(b, out); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}
