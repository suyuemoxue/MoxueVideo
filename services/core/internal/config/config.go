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
	Env           string
	HTTPAddr      string
	MySQLDSN      string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RabbitMQURL   string
	ChatGRPCAddr  string
	OSS           OSSConfig
}

type OSSConfig struct {
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	RoleARN         string `yaml:"roleArn"`
	RoleSessionName string `yaml:"roleSessionName"`
	DurationSeconds int    `yaml:"durationSeconds"`
}

func Load() Config {
	cfg := Config{
		Env:           getString("APP_ENV", "dev"),
		HTTPAddr:      getString("HTTP_ADDR", ":8080"),
		MySQLDSN:      getString("MYSQL_DSN", "moxue:suyuemoxue-mojianxue@tcp(127.0.0.1:3307)/moxuevideo?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:     getString("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getString("REDIS_PASSWORD", "suyuemoxue-mojianxue"),
		RedisDB:       getInt("REDIS_DB", 0),
		RabbitMQURL:   getString("RABBITMQ_URL", "amqp://app:apppass@127.0.0.1:5672/"),
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
		return fmt.Errorf("read config file: %w", err)
	}
	if err := yaml.Unmarshal(b, out); err != nil {
		return fmt.Errorf("unmarshal config yaml: %w", err)
	}
	return nil
}
