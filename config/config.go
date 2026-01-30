package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		HTTP    HTTP
		Log     Log
		PG      PG
		Redis   Redis
		Swagger Swagger
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		URL     string `env:"PG_URL,required"`
		PoolMax int    `env:"PG_POOL_MAX,required"`
	}

	Redis struct {
		Addr        string `env:"REDIS_ADDR,required"`
		DB          int    `env:"REDIS_DB,required"`
		User        string `env:"REDIS_USER"`
		Password    string `env:"REDIS_PASSWORD"`
		DialTimeout int    `env:"REDIS_DIAL_TIMEOUT"`
		Timeout     int    `env:"REDIS_TIMEOUT"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %v", err)
	}

	return cfg, nil
}
