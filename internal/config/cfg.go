package config

import (
	"flag"
	"sync"
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

var (
	cfg  *Config
	once sync.Once
)

func NewConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
		flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")
		flag.Parse()
	})

	return cfg
}
