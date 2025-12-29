package config

import (
	"encoding/json"
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
	"os"
)

// Config is a struct that holds the configuration for the application.
type Config struct {
	ServerAddr      string `env:"SERVER_ADDRESS" json:"server_address" envDefault:"localhost:8080"`
	GRPCAddr        string `env:"GRPC_ADDRESS" json:"grpc_address" envDefault:"localhost:8081"`
	BaseURL         string `env:"BASE_URL" json:"base_url" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DSN             string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	Config          string `env:"CONFIG"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"abcdefghijklmnopqrstuvwxyz123456"`
	AuditFile       string `env:"AUDIT_FILE"`
	AuditURL        string `env:"AUDIT_URL"`
	CertFile        string `env:"CERT_FILE" envDefault:"certs/cert.pem"`
	KeyFile         string `env:"KEY_FILE" envDefault:"certs/key.pem"`
	Subnet          string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

// NewConfig creates a new Config object.
func NewConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Printf("failed to load envs: %s", err.Error())
	}

	flag.StringVar(&cfg.Config, "c", cfg.Config, "Config file path")
	flag.StringVar(&cfg.Config, "config", cfg.Config, "Config file path")
	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server address")
	flag.StringVar(&cfg.GRPCAddr, "g", cfg.GRPCAddr, "GRPC server address")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links")
	flag.StringVar(&cfg.DSN, "d", cfg.DSN, "DSN for database connection")
	flag.StringVar(&cfg.SecretKey, "k", cfg.SecretKey, "Secret key for cryptographic")
	flag.StringVar(&cfg.AuditFile, "audit-file", cfg.AuditFile, "Audit file path")
	flag.StringVar(&cfg.AuditURL, "audit-url", cfg.AuditURL, "Audit URL")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "Enable HTTPS")
	flag.StringVar(&cfg.CertFile, "cert-file", cfg.CertFile, "Certificate file path")
	flag.StringVar(&cfg.KeyFile, "key-file", cfg.KeyFile, "Key file path")
	flag.StringVar(&cfg.Subnet, "t", cfg.Subnet, "Subnet")
	flag.Parse()

	if cfg.Config != "" {
		data, err := os.ReadFile(cfg.Config)
		if err != nil {
			log.Printf("failed to read config file: %s", err.Error())
		} else {
			if err := json.Unmarshal(data, cfg); err != nil {
				log.Printf("failed to unmarshal config file: %s", err.Error())
			}
		}
	}

	log.Printf("CFG: %v\n", cfg)
	return cfg
}
