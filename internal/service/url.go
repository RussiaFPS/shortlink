package service

import (
	"fmt"
	"github.com/RussiaFPS/shortlink/internal/config"
	"github.com/RussiaFPS/shortlink/internal/model"
	"github.com/RussiaFPS/shortlink/internal/repository"
	"log"
	"math/rand"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type IURLShortenerService interface {
	Shorten(originalURL string) (string, error)
	GetOriginal(shortURL string) (string, bool)
	ShorterMulti(req []model.MultiReq) ([]model.MultiResp, error)
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
	log.Printf("URLShortenerService:Shorten with originalURL: %s", originalURL)

	if shortURL, ok := s.r.GetShortURL(originalURL); ok {
		return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL), nil
	}

	URL, err := s.r.StorageURL(originalURL, s.generateShortID())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, URL), nil
}

func (s *URLShortenerService) GetOriginal(shortID string) (string, bool) {
	log.Printf("URLShortenerService:GetOriginal with shortID: %s", shortID)
	return s.r.GetOriginalURL(shortID)
}

func (s *URLShortenerService) ShorterMulti(req []model.MultiReq) ([]model.MultiResp, error) {
	data := make([]model.URLStorage, 0)

	for _, v := range req {
		d := model.URLStorage{
			UUID:        v.CorrID,
			OriginalURL: v.URL,
		}

		if shortURL, ok := s.r.GetShortURL(v.URL); ok {
			d.ShortURL = shortURL
		} else {
			d.ShortURL = s.generateShortID()
		}

		data = append(data, d)
	}

	resp, err := s.r.StoreMultiURL(data)
	if err != nil {
		return nil, err
	}

	for i := range resp {
		if !strings.Contains(resp[i].ShortURL, "http") {
			resp[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, resp[i].ShortURL)
		}
	}

	return resp, nil
}
