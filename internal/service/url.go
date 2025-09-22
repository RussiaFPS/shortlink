package service

import (
	"context"
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
	Shorten(ctx context.Context, originalURL string, userID string) (string, bool, error)
	GetOriginal(ctx context.Context, shortURL string) (*model.URLStorage, bool)
	ShorterMulti(ctx context.Context, req []model.MultiReq, userID string) ([]model.MultiResp, error)
	PingDB(ctx context.Context) error
	randShortID() string
	generateShortID(ctx context.Context) string
	GetShortenedURLByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error)
	DeleteURLs(req *[]string, userID string) error
}

type URLShortenerService struct {
	cfg      *config.Config
	shortLen int
	r        repository.IURLShortenerRepository
}

func NewURLShortener(ctx context.Context, cfg *config.Config, shortLen int) IURLShortenerService {
	return &URLShortenerService{
		shortLen: shortLen,
		cfg:      cfg,
		r:        repository.New(ctx, cfg),
	}
}

func (s *URLShortenerService) PingDB(ctx context.Context) error {
	return s.r.Ping(ctx)
}

func (s *URLShortenerService) generateShortID(ctx context.Context) string {
	var id string
	for {
		id = s.randShortID()
		if _, exists := s.r.GetOriginalURL(ctx, id); !exists {
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

func (s *URLShortenerService) Shorten(ctx context.Context, originalURL string, userID string) (string, bool, error) {
	log.Printf("URLShortenerService:Shorten with originalURL: %s", originalURL)

	if shortURL, ok := s.r.GetShortURL(ctx, originalURL); ok {
		return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortURL), true, nil
	}

	URL, err := s.r.StorageURL(ctx, originalURL, s.generateShortID(ctx), userID)
	if err != nil {
		return "", false, err
	}

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, URL), false, nil
}

func (s *URLShortenerService) GetOriginal(ctx context.Context, shortID string) (*model.URLStorage, bool) {
	log.Printf("URLShortenerService:GetOriginal with shortID: %s", shortID)
	return s.r.GetOriginalURL(ctx, shortID)
}

func (s *URLShortenerService) ShorterMulti(ctx context.Context, req []model.MultiReq, userID string) ([]model.MultiResp, error) {
	data, oldData := make([]model.URLStorage, 0), make([]model.MultiResp, 0)

	for _, v := range req {
		d := model.URLStorage{
			UUID:        v.CorrID,
			OriginalURL: v.URL,
			UserID:      userID,
		}

		if shortURL, ok := s.r.GetShortURL(ctx, v.URL); ok {
			oldData = append(oldData, model.MultiResp{CorrID: v.CorrID, ShortURL: shortURL})
		} else {
			d.ShortURL = s.generateShortID(ctx)
			data = append(data, d)
		}
	}

	resp, err := s.r.StoreMultiURL(ctx, data)
	if err != nil {
		return nil, err
	}

	resp = append(resp, oldData...)
	for i := range resp {
		if !strings.Contains(resp[i].ShortURL, "http") {
			resp[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, resp[i].ShortURL)
		}
	}

	return resp, nil
}

func (s *URLShortenerService) GetShortenedURLByUserID(ctx context.Context, userID string) ([]model.RespUserURL, error) {
	data, err := s.r.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range data {
		if !strings.Contains(data[i].ShortURL, "http") {
			data[i].ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, data[i].ShortURL)
		}
	}
	return data, nil
}

func (s *URLShortenerService) DeleteURLs(req *[]string, userID string) error {
	return s.r.DeleteRecords(*req, userID)
}
