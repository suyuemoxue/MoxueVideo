package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env           string    `yaml:"env"`
	HTTPAddr      string    `yaml:"http_addr"`
	MySQLDSN      string    `yaml:"mysql_dsn"`
	RedisAddr     string    `yaml:"redis_addr"`
	RedisPassword string    `yaml:"redis_password"`
	RedisDB       int       `yaml:"redis_db"`
	RabbitMQURL   string    `yaml:"rabbitmq_url"`
	ChatGRPCAddr  string    `yaml:"chat_grpc_addr"`
	OSS           OSSConfig `yaml:"oss"`
}

type OSSConfig struct {
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	RoleARN         string `yaml:"role_arn"`
	RoleSessionName string `yaml:"role_session_name"`
	DurationSeconds int    `yaml:"duration_seconds"`
}

func Load() Config {
	cfg := Config{
		Env:           getString("APP_ENV", "dev"),
		HTTPAddr:      getString("HTTP_ADDR", ":8080"),
		MySQLDSN:      getString("MYSQL_DSN", ""),
		RedisAddr:     getString("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getString("REDIS_PASSWORD", ""),
		RedisDB:       getInt("REDIS_DB", 0),
		RabbitMQURL:   getString("RABBITMQ_URL", ""),
		ChatGRPCAddr:  getString("CHAT_GRPC_ADDR", "127.0.0.1:50051"),
	}

	if path := getString("CONFIG_FILE", ""); path != "" {
		_ = loadFromYAML(path, &cfg)
	} else {
		_ = loadFromYAML("config.local.yaml", &cfg)
	}

	cfg.Env = getString("APP_ENV", cfg.Env)
	cfg.HTTPAddr = getString("HTTP_ADDR", cfg.HTTPAddr)
	cfg.MySQLDSN = getString("MYSQL_DSN", cfg.MySQLDSN)
	cfg.RedisAddr = getString("REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisPassword = getString("REDIS_PASSWORD", cfg.RedisPassword)
	cfg.RedisDB = getInt("REDIS_DB", cfg.RedisDB)
	cfg.RabbitMQURL = getString("RABBITMQ_URL", cfg.RabbitMQURL)
	cfg.ChatGRPCAddr = getString("CHAT_GRPC_ADDR", cfg.ChatGRPCAddr)
	cfg.OSS.Region = getString("OSS_REGION", cfg.OSS.Region)
	cfg.OSS.Endpoint = getString("OSS_ENDPOINT", cfg.OSS.Endpoint)
	cfg.OSS.Bucket = getString("OSS_BUCKET", cfg.OSS.Bucket)
	cfg.OSS.AccessKeyID = getString("OSS_ACCESS_KEY_ID", cfg.OSS.AccessKeyID)
	cfg.OSS.AccessKeySecret = getString("OSS_ACCESS_KEY_SECRET", cfg.OSS.AccessKeySecret)
	cfg.OSS.RoleARN = getString("OSS_ROLE_ARN", cfg.OSS.RoleARN)
	cfg.OSS.RoleSessionName = getString("OSS_ROLE_SESSION_NAME", cfg.OSS.RoleSessionName)
	cfg.OSS.DurationSeconds = getInt("OSS_STS_DURATION_SECONDS", cfg.OSS.DurationSeconds)

	if cfg.OSS.Region == "" {
		cfg.OSS.Region = "cn-hangzhou"
	}
	if cfg.OSS.RoleSessionName == "" {
		cfg.OSS.RoleSessionName = "moxuevideo-core"
	}
	if cfg.OSS.DurationSeconds == 0 {
		cfg.OSS.DurationSeconds = 900
	}

	return cfg
}

func getString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
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
