package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerAddr string
	BaseURL    string
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("config", flag.ContinueOnError)

	fs.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
	fs.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for shortened links")

	if err := fs.Parse(os.Args); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}
	return cfg, nil
}
