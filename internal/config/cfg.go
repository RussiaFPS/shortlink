package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

// Config is a struct that holds the configuration for the application.
type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"URL.json"`
	DSN             string `env:"DATABASE_DSN" envDefault:"postgres://postgres:qwer1234@localhost:5432/shortlink"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"abcdefghijklmnopqrstuvwxyz123456"`
	AuditFile       string `env:"AUDIT_FILE"`
	AuditURL        string `env:"AUDIT_URL"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
	CertFile        string `env:"CERT_FILE" envDefault:"certs/cert.pem"`
	KeyFile         string `env:"KEY_FILE" envDefault:"certs/key.pem"`
}

// NewConfig creates a new Config object.
func NewConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Printf("failed to load envs: %s", err.Error())
	}

	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server address")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links")
	flag.StringVar(&cfg.DSN, "d", cfg.DSN, "DSN for database connection")
	flag.StringVar(&cfg.SecretKey, "k", cfg.SecretKey, "Secret key for cryptographic")
	flag.StringVar(&cfg.AuditFile, "audit-file", cfg.AuditFile, "Audit file path")
	flag.StringVar(&cfg.AuditURL, "audit-url", cfg.AuditURL, "Audit URL")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "Enable HTTPS")
	flag.StringVar(&cfg.CertFile, "cert-file", cfg.CertFile, "Certificate file path")
	flag.StringVar(&cfg.KeyFile, "key-file", cfg.KeyFile, "Key file path")
	flag.Parse()

	log.Printf("CFG: %v\n", cfg)
	return cfg
}
