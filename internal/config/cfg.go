package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"URL.json"`
	DSN             string `env:"DATABASE_DSN" envDefault:"postgres://postgres:qwer1234@localhost:5432/shortlink"`
}

func NewConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Printf("failed to load envs: %s", err.Error())
	}

	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server address")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links")
	flag.StringVar(&cfg.DSN, "d", cfg.DSN, "DSN for database connection")
	flag.Parse()

	log.Printf("CFG: %v\n", cfg)
	return cfg
}
