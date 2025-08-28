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
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
		flag.StringVar(&cfg.FileStoragePath, "f", "URL.json", "File storage path")
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")

		flag.Parse()

		err := env.Parse(cfg)
		if err != nil {
			log.Printf("failed to load envs: %s", err.Error())
		}

		log.Printf("CFG: %v\n", cfg)
	})

	return cfg
}
