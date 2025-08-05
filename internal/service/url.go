package service

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLShortenerService struct {
	cfg      *config.Config
	shortLen int
	r        *repository.URLShortenerRepository
}

func NewURLShortener(cfg *config.Config, shortLen int) *URLShortenerService {
	return &URLShortenerService{
		shortLen: shortLen,
		cfg:      cfg,
		r:        repository.NewURLShortenerRepository(cfg),
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
	if shortURL, ok := s.r.GetShortURL(originalURL); ok {
		return shortURL, nil
	}

	return s.r.StorageURL(originalURL, s.generateShortID()), nil
}

func (s *URLShortenerService) GetOriginal(shortID string) (string, bool) {
	return s.r.GetOriginalURL(shortID)
}
