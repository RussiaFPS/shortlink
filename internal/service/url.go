package service

import (
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"math/rand"
	"net/url"
	"sync"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLShortenerService struct {
	cfg      *config.Config
	mu       sync.Mutex
	urls     map[string]string // {shortID: originalURL}
	shortLen int
}

func NewURLShortener(cfg *config.Config, shortLen int) *URLShortenerService {
	return &URLShortenerService{
		urls:     make(map[string]string),
		shortLen: shortLen,
		cfg:      cfg,
	}
}

func (s *URLShortenerService) generateShortID() string {
	id := make([]byte, s.shortLen)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)
}

func (s *URLShortenerService) Shorten(originalURL string) (string, error) {
	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return "", fmt.Errorf("invalid URL")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for shortID, existingURL := range s.urls {
		if existingURL == originalURL {
			return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortID), nil
		}
	}

	shortID := s.generateShortID()
	s.urls[shortID] = originalURL

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortID), nil
}

func (s *URLShortenerService) GetOriginal(shortID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalURL, exists := s.urls[shortID]
	return originalURL, exists
}
