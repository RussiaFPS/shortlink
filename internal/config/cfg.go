package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
	"sync"
)

type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := env.Parse(cfg); err != nil {
			log.Fatal(err)
		}

		if cfg.FileStoragePath == "" {
			flag.StringVar(&cfg.FileStoragePath, "f", "URL.json", "File storage path")
		}
		if cfg.ServerAddr == "" {
			flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
		}
		if cfg.BaseURL == "" {
			flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
		}

		flag.Parse()
	})

	return cfg
}
