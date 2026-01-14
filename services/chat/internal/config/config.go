package config

import "os"

type Config struct {
	Env         string
	GRPCAddr    string
	WSAddr      string
	MySQLDSN    string
	RabbitMQURL string
}

func Load() Config {
	cfg := Config{
		Env:         getString("APP_ENV", "dev"),
		GRPCAddr:    getString("GRPC_ADDR", ":50051"),
		WSAddr:      getString("WS_ADDR", ":50052"),
		MySQLDSN:    getString("MYSQL_DSN", ""),
		RabbitMQURL: getString("RABBITMQ_URL", ""),
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
