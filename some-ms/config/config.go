package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `yaml:"level"`
	}

	PG struct {
		URL string `env:"PG_URL" env-required:"true"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
