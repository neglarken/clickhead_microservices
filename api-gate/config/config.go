package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		Services
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

	Services struct {
		SomeMsServerAddress string
		AuthMsServerAddress string
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, err
	}

	cfg.SomeMsServerAddress = fmt.Sprintf("%s:%s", os.Getenv("SOME_MS_HOST"), os.Getenv("SOME_MS_PORT"))
	cfg.AuthMsServerAddress = fmt.Sprintf("%s:%s", os.Getenv("AUTH_MS_HOST"), os.Getenv("AUTH_MS_PORT"))

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
