package repository

import (
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"sync"
)

type URLShortenerRepository struct {
	cfg        *config.Config
	mu         sync.RWMutex
	urls       map[string]string // {shortID: originalURL}
	urlToShort map[string]string // {originalURL: shortID}
}

func NewURLShortenerRepository(cfg *config.Config) *URLShortenerRepository {
	return &URLShortenerRepository{
		urls:       make(map[string]string),
		urlToShort: make(map[string]string),
		cfg:        cfg,
	}
}

func (r *URLShortenerRepository) GetShortURL(originalURL string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.urlToShort[originalURL]; ok {
		return fmt.Sprintf("%s/%s", r.cfg.BaseURL, v), true
	}

	return "", false
}

func (r *URLShortenerRepository) StorageURL(originalURL string, shortID string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.urls[shortID] = originalURL
	r.urlToShort[originalURL] = shortID

	return fmt.Sprintf("%s/%s", r.cfg.BaseURL, shortID)
}

func (r *URLShortenerRepository) GetOriginalURL(shortID string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.urls[shortID]; ok {
		return v, true
	}
	return "", false
}
