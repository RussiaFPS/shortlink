package service

import (
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type IURLShortenerService interface {
	Shorten(originalURL string) (string, error)
	GetOriginal(shortURL string) (string, bool)
	PingDB() error
	randShortID() string
	generateShortID() string
}

type URLShortenerService struct {
	cfg      *config.Config
	shortLen int
	r        repository.IURLShortenerRepository
}

func NewURLShortener(cfg *config.Config, shortLen int) IURLShortenerService {
	return &URLShortenerService{
		shortLen: shortLen,
		cfg:      cfg,
		r:        repository.New(cfg),
	}
}

func (s *URLShortenerService) PingDB() error {
	return s.r.Ping()
}

func (s *URLShortenerService) generateShortID() string {
	var id string
	for {
		id = s.randShortID()
		if _, exists := s.r.GetOriginalURL(id); !exists {
			break
		}
	}
	return id
}

func (s *URLShortenerService) randShortID() string {
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

	URL, err := s.r.StorageURL(originalURL, s.generateShortID())
	if err != nil {
		return "", err
	}

	return URL, nil
}

func (s *URLShortenerService) GetOriginal(shortID string) (string, bool) {
	return s.r.GetOriginalURL(shortID)
}
