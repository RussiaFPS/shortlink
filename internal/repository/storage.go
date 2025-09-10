package repository

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository/database"
	"github.com/RussiaFPS/shortlink/internal/repository/memory"
	"log"
)

type IURLShortenerRepository interface {
	GetShortURL(originalURL string) (string, bool)
	StorageURL(originalURL string, shortID string) (string, error)
	GetOriginalURL(shortID string) (string, bool)
	StoreMultiURL(req []model.URLStorage) ([]model.MultiResp, error)
	Ping() error
}

func New(cfg *config.Config) IURLShortenerRepository {
	var storage IURLShortenerRepository
	var err error

	if cfg.DSN != "" {
		storage, err = database.New(cfg)
		if err != nil {
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
