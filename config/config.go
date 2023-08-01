package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath string = "config/.env"

type (
	Config struct {
		MySql
		PG
		RMQ
		HTTP
		WorkerPool
	}

	MySql struct {
		Host     string `env-required:"true" env:"MY_HOST"`
		Username string `env-required:"true" env:"MY_USERNAME"`
		Password string `env-required:"true" env:"MY_PASSWORD"`
		Port     string `env-required:"true" env:"MY_PORT"      env-default:"3306"`
		DBName   string `env-required:"true" env:"MY_DBNAME"`
	}

	PG struct {
		Host        string `env-required:"true" env:"PG_HOST"`
		Username    string `env-required:"true" env:"PG_USERNAME"`
		Password    string `env-required:"true" env:"PG_PASSWORD"`
		Port        string `env-required:"true" env:"PG_PORT"      env-default:"5432"`
		DBName      string `env-required:"true" env:"PG_DBNAME"`
		MaxOpenConn int32  `env-required:"true" env:"PG_MAXOPENCONN"`
	}

	RMQ struct {
		Host     string `env-required:"true" env:"RMQ_HOST"`
		Username string `env-required:"true" env:"RMQ_USERNAME"`
		Password string `env-required:"true" env:"RMQ_PASSWORD"`
		Port     string `env-required:"true" env:"RMQ_PORT"      env-default:"5672"`
	}

	HTTP struct {
		HTTPAddr string `env-required:"true" env:"HTTP_ADDR"`
	}

	WorkerPool struct {
		WorkerCount int `env-required:"true" env:"WORKER_COUNT"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
