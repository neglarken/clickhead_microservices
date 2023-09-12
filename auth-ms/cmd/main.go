package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/neglarken/clickhead/auth-ms/config"
	"github.com/neglarken/clickhead/auth-ms/internal/httpserver"
)

var (
	configPath string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	flag.StringVar(&configPath, "config-path", "config/config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	config, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := httpserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
