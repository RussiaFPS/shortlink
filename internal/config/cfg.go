package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"sync"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS,required"`
	BaseURL    string `env:"BASE_URL,required"`
	//FileStoragePath string `env:"FILE_STORAGE_PATH,required"`
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		if err := env.Parse(cfg); err == nil {
			return
		}

		if cfg.ServerAddr == "" {
			flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
		}
		if cfg.BaseURL == "" {
			flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
		}
		//if cfg.FileStoragePath == "" {
		//	flag.StringVar(&cfg.FileStoragePath, "f", "storage.json", "File storage path")
		//}

		flag.Parse()
	})

	return cfg
}
