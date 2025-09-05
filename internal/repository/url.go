package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/config/db/postgres"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
	"time"
)

type IURLShortenerRepository interface {
	GetShortURL(originalURL string) (string, bool)
	StorageURL(originalURL string, shortID string) (string, error)
	GetOriginalURL(shortID string) (string, bool)
	PingDB() error
	loadRecords() error
}

type URLShortenerRepository struct {
	cfg        *config.Config
	db         *pgxpool.Pool
	mu         sync.RWMutex
	urls       map[string]string // {shortID: originalURL}
	urlToShort map[string]string // {originalURL: shortID}
	records    []model.URLStorage
	storage    IStorage
}

func NewURLShortenerRepository(cfg *config.Config) IURLShortenerRepository {
	s := &URLShortenerRepository{
		urls:       make(map[string]string),
		urlToShort: make(map[string]string),
		records:    make([]model.URLStorage, 0),
		cfg:        cfg,
		storage:    NewFileStorage(cfg),
	}

	dbConn, err := postgres.NewPostgres(context.Background(), cfg.DSN)
	if err != nil {
		log.Fatalf("failed to initialize dbConn: %v", err)
	}
	s.db = dbConn

	if err = s.loadRecords(); err != nil {
		log.Fatalf("failed to load records: %v", err)
	}

	return s
}

func (r *URLShortenerRepository) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	return r.db.Ping(ctx)
}

func (r *URLShortenerRepository) GetShortURL(originalURL string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.urlToShort[originalURL]; ok {
		return fmt.Sprintf("%s/%s", r.cfg.BaseURL, v), true
	}

	return "", false
}

func (r *URLShortenerRepository) StorageURL(originalURL string, shortID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s := model.URLStorage{
		UUID:        shortID,
		ShortURL:    shortID,
		OriginalURL: originalURL,
	}

	r.records = append(r.records, s)
	newData, err := json.MarshalIndent(r.records, "", "   ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %v", err)
	}

	if err = r.storage.SaveDataToFile(newData); err != nil {
		return "", fmt.Errorf("failed to save data to file: %v", err)
	}

	r.urls[shortID] = originalURL
	r.urlToShort[originalURL] = shortID

	return fmt.Sprintf("%s/%s", r.cfg.BaseURL, shortID), nil
}

func (r *URLShortenerRepository) GetOriginalURL(shortID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.urls[shortID]; ok {
		return v, true
	}
	return "", false
}

func (r *URLShortenerRepository) loadRecords() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.storage.GetDataFromFile()
	if err != nil {
		return fmt.Errorf("failed to load data from file: %v", err)
	}

	if len(data) != 0 {
		if err = json.Unmarshal(data, &r.records); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	for _, record := range r.records {
		r.urls[record.ShortURL] = record.OriginalURL
		r.urlToShort[record.OriginalURL] = record.ShortURL
	}

	log.Printf("Loaded %d URLs from storage file", len(r.records))
	return nil
}
