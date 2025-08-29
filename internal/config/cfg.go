package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
	"sync"
)

type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"URL.json"`
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := env.Parse(cfg); err != nil {
			log.Printf("failed to load envs: %s", err.Error())
		}

		flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server address")
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links")
		flag.Parse()

		log.Printf("CFG: %v\n", cfg)
	})

	return cfg
}
