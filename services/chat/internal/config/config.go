package config

import "os"

type Config struct {
	Env      string
	GRPCAddr string
}

func Load() Config {
	cfg := Config{
		Env:      getString("APP_ENV", "dev"),
		GRPCAddr: getString("GRPC_ADDR", ":50051"),
	}

	cfg.Env = getString("APP_ENV", cfg.Env)
	cfg.GRPCAddr = getString("GRPC_ADDR", cfg.GRPCAddr)
	return cfg
}

func getString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
