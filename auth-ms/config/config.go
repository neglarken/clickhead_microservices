package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
		JWT
		Hasher
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
		URL string `yaml:"url" env-required:"true" env:"PG_URL"`
	}

	JWT struct {
		SignKey         string        `env-required:"true" env:"JWT_SIGN_KEY"`
		AccessTokenTTL  time.Duration `env-required:"true" yaml:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `env-required:"true" yaml:"refreshTokenTTL"`
	}

	Hasher struct {
		Salt string `env-required:"true" env:"HASHER_SALT"`
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
