package repository

import (
	"context"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository/database"
	"github.com/RussiaFPS/shortlink/internal/repository/memory"
	"log"
)

type IURLShortenerRepository interface {
	GetShortURL(ctx context.Context, originalURL string) (string, bool)
	StorageURL(ctx context.Context, originalURL string, shortID string) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (string, bool)
	StoreMultiURL(ctx context.Context, req []model.URLStorage) ([]model.MultiResp, error)
	Ping(ctx context.Context) error
}

func New(ctx context.Context, cfg *config.Config) IURLShortenerRepository {
	var storage IURLShortenerRepository
	var err error

	if cfg.DSN != "" {
		storage, err = database.New(ctx, cfg)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			storage, err = memory.New(cfg)
			if err != nil {
				log.Fatalf("failed init memory storage: %v", err)
			}
		}
	} else {
		storage, err = memory.New(cfg)
		if err != nil {
			log.Fatalf("failed init memory storage: %v", err)
		}
	}

	return storage
}
