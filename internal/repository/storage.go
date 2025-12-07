package repository

import (
	"context"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository/database"
	"github.com/RussiaFPS/shortlink/internal/repository/memory"
	"log"
)

// URLRepository is an interface for URL storage.
type URLRepository interface {
	GetShortURL(ctx context.Context, originalURL string) (string, bool)
	StorageURL(ctx context.Context, originalURL string, shortID string, userID string) (string, error)
	GetOriginalURL(ctx context.Context, shortID string) (*model.URLStorage, bool)
	StoreMultiURL(ctx context.Context, req []model.URLStorage) ([]model.MultiResp, error)
	Ping(ctx context.Context) error
	FindAllByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error)
	DeleteRecords(ids []string, userID string) error
	Close()
}

// New creates a new URLRepository based on the provided configuration.
func New(ctx context.Context, cfg *config.Config) URLRepository {
	var storage URLRepository
	var err error

	if cfg.DSN != "" {
		ch := make(chan model.DellURL)
		storage, err = database.New(ctx, cfg, ch)
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
