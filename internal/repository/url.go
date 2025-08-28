package repository

import (
	"encoding/json"
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"log"
	"os"
	"sync"
)

type URLShortenerRepository struct {
	cfg        *config.Config
	mu         sync.RWMutex
	urls       map[string]string // {shortID: originalURL}
	urlToShort map[string]string // {originalURL: shortID}
	file       *os.File
	records    []model.URLStorage
}

func NewURLShortenerRepository(cfg *config.Config) *URLShortenerRepository {
	s := &URLShortenerRepository{
		urls:       make(map[string]string),
		urlToShort: make(map[string]string),
		records:    make([]model.URLStorage, 0),
		cfg:        cfg,
	}

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open storage file: %v", err)
	}
	s.file = file

	if err = s.loadFromFile(); err != nil {
		log.Fatalf("failed to load data from file: %v", err)
	}

	return s
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

	newData, err := json.MarshalIndent(r.records, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	if err = os.WriteFile(r.cfg.FileStoragePath, newData, 0644); err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
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

func (r *URLShortenerRepository) loadFromFile() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile(r.cfg.FileStoragePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	if err = json.Unmarshal(data, &r.records); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, record := range r.records {
		r.urls[record.ShortURL] = record.OriginalURL
		r.urlToShort[record.OriginalURL] = record.ShortURL
	}

	log.Printf("Loaded %d URLs from storage file", len(r.records))
	return nil
}
